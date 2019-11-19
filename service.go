package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

// Defaults that can be overridden.
const DefaultFlushInterval = 15 * time.Second
const DefaultWorkers = 1
const DefaultBulkActions = 1000
const DefaultMaxRetries = 5
const DefaultElasticsearchURL = "http://elasticsearch:9200"
const DefaultRunOnce = false

type Service struct {
	processor *elastic.BulkProcessor
}

type Doc struct {
	Message string                  `json:"message"`
	Level   string                  `json:"level"`
	Date    string                  `json:"date"`
	Data    *map[string]interface{} `json:"data"`
}

func (svc *Service) Init() error {
	var err error
	var flushInterval time.Duration
	var workers int
	var bulkActions int
	var maxRetries int
	var elasticsearchUrl string
	var runOnce bool

	// Check if there is a FLUSH_INTERVAL passed in.
	if os.Getenv("FLUSH_INTERVAL") == "" {
		flushInterval = DefaultFlushInterval
	} else {
		i, err := strconv.Atoi(os.Getenv("FLUSH_INTERVAL"))
		if err != nil {
			return err
		}

		flushInterval = time.Duration(i) * time.Second
	}

	// Check if there is a WORKERS passed in.
	if os.Getenv("WORKERS") == "" {
		workers = DefaultWorkers
	} else {
		workers, err = strconv.Atoi(os.Getenv("WORKERS"))
		if err != nil {
			return err
		}
	}

	// Check if there is a BULK_ACTIONS passed in.
	if os.Getenv("BULK_ACTIONS") == "" {
		bulkActions = DefaultBulkActions
	} else {
		bulkActions, err = strconv.Atoi(os.Getenv("BULK_ACTIONS"))
		if err != nil {
			return err
		}
	}

	// Check for the MAX_RETRIES
	if os.Getenv("MAX_RETRIES") == "" {
		maxRetries = DefaultMaxRetries
	} else {
		maxRetries, err = strconv.Atoi(os.Getenv("MAX_RETRIES"))
		if err != nil {
			return err
		}
	}

	// Check for the ELASTICSEARCH_URL
	if os.Getenv("ELASTICSEARCH_URL") == "" {
		elasticsearchUrl = DefaultElasticsearchURL
	} else {
		elasticsearchUrl = os.Getenv("ELASTICSEARCH_URL")
	}

	// Check for the RUN_ONCE
	if os.Getenv("RUN_ONCE") == "" {
		runOnce = DefaultRunOnce
	} else {
		runOnce, err = strconv.ParseBool(os.Getenv("RUN_ONCE"))
		if err != nil {
			return err
		}
	}

	// Setup the elastic client.
	client, err := elastic.NewClient(
		elastic.SetURL(elasticsearchUrl),
		elastic.SetSniff(false),
		elastic.SetMaxRetries(maxRetries),
	)
	if err != nil {
		return err
	}

	// Setup the bulk processor to flush the logs to elastic.
	processor, err := client.BulkProcessor().
		Name("processor").
		Workers(workers).
		BulkActions(bulkActions).
		BulkSize(2 << 20).
		FlushInterval(flushInterval).
		Do(context.Background())
	if err != nil {
		return err
	}
	defer processor.Close()

	// Setup the service processor.
	svc.processor = processor

	// Setup http server.
	mux := http.NewServeMux()

	// Setup the http handle for incoming logs.
	mux.HandleFunc("/", svc.handler)
	mux.HandleFunc("/ping", svc.ping)

	// Run once mode for unit tests.
	if runOnce {
		return nil
	}

	return http.ListenAndServe(":3000", mux)
}

// Private

func (svc *Service) ping(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
}

func (svc *Service) handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method == "POST" {
		// Decode the incoming payload that will be placed into the bulker.
		decoder := json.NewDecoder(r.Body)

		var doc Doc
		err := decoder.Decode(&doc)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Break up the path to get the index.
		path := strings.Split(r.URL.Path, "/")

		// Must include a path which will be the index and type you are targeting.
		if len(path) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t := time.Now()

		doc.Date = t.Format(time.RFC3339)

		r := elastic.NewBulkIndexRequest().Index(fmt.Sprintf(`%s-%s`, path[1], t.Format("2006-01-02"))).Type(path[2]).Doc(doc)
		svc.processor.Add(r)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

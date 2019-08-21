package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

func TestHttpServer(t *testing.T) {
	client, _ := elastic.NewClient()

	svc := &Service{}
	processor, _ := client.BulkProcessor().Name("processor").Do(context.Background())
	svc.processor = processor

	payload := `{ "message" : "test message", "level" : "error" }`
	body := strings.NewReader(payload)

	req := httptest.NewRequest("POST", "/logs/misc", body)
	w := httptest.NewRecorder()
	svc.handler(w, req)

	resp := w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusNoContent)

	req = httptest.NewRequest("POST", "/", body)
	w = httptest.NewRecorder()
	svc.handler(w, req)

	resp = w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	req = httptest.NewRequest("GET", "/logs/misc", nil)
	w = httptest.NewRecorder()
	svc.handler(w, req)

	resp = w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusMethodNotAllowed)

	payload = `fail`
	body = strings.NewReader(payload)

	req = httptest.NewRequest("POST", "/logs/misc", nil)
	w = httptest.NewRecorder()
	svc.handler(w, req)

	resp = w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)

	payload = `{ "message" : "test message", "level" : "error" }`
	body = strings.NewReader(payload)

	req = httptest.NewRequest("POST", "/logs/misc/bad/url", body)
	w = httptest.NewRecorder()
	svc.handler(w, req)

	resp = w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusBadRequest)
}

func TestInit(t *testing.T) {
	svc := Service{}

	err := svc.Init()
	assert.Equal(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "1")
	os.Setenv("WORKERS", "1")
	os.Setenv("BULK_ACTIONS", "1")
	os.Setenv("MAX_RETRIES", "5")
	os.Setenv("ELASTICSEARCH_URL", "http://127.0.0.1:9200")
	os.Setenv("RUN_ONCE", "true")

	err = svc.Init()
	assert.Equal(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "1")
	os.Setenv("WORKERS", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "1")
	os.Setenv("WORKERS", "1")
	os.Setenv("BULK_ACTIONS", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "1")
	os.Setenv("WORKERS", "1")
	os.Setenv("BULK_ACTIONS", "1")
	os.Setenv("MAX_RETRIES", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)

	os.Setenv("FLUSH_INTERVAL", "1")
	os.Setenv("WORKERS", "1")
	os.Setenv("BULK_ACTIONS", "1")
	os.Setenv("MAX_RETRIES", "1")
	os.Setenv("RUN_ONCE", "fail")

	err = svc.Init()
	assert.NotEqual(t, err, nil)
}

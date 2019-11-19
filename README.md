# Logger

Simple lightweight API that accepts a POST request with a JSON payload that is
written to a buffer and flushed at configurable intervals into Elasticsearch.
Whichever setting (time, buffer size) that comes first for the processor
configuration will execute the flush. You can start the entire service with
`docker-compose up`.

[![CircleCI](https://circleci.com/gh/cachelab/logger.svg?style=svg)](https://circleci.com/gh/cachelab/logger)

## Usage

This task is configured by the following environment variables:

```bash
FLUSH_INTERVAL    # the number of seconds you want the buffer to be flushed
WORKERS           # the number of workers that will be available to buffer requests
BULK_ACTIONS      # the maximum number of records that will be flushed
MAX_RETRIES       # how many times the client will try to connect to Elasticsearch
ELASTICSEARCH_URL # the url to the Elasticsearch cluster
RUN_ONCE          # used for unit testing to not start the http server
```

## Example

```
curl -XPOST http://127.0.0.1:3000/your-api/log -d '{"message": "Error in line 128. Bad object.", "level": "error", "data": {"http_status": 500}}}'
```

![alt text](/images/screenshot.png)

## Contributing

* `make run` - runs the api in a docker container
* `make build` - builds your logger docker container
* `make vet` - go fmt and vet code
* `make test` - run unit tests

Before you submit a pull request please update the semantic version inside of
`main.go` with what you feel is appropriate and then edit the `CHANGELOG.md` with
your changes and follow a similar structure to what is there.

NAME := logger
MAINTAINER := cachelab
VERSION := $(shell grep "const version =" main.go | cut -d\" -f2)

.PHONY: *

default: build

up:
	@echo Running Compose
	docker-compose up

run: build
	@echo Running Container
	docker run -e "ELASTICSEARCH_URL=http://elasticsearch:9200" --network logger_logger -p 3000:3000 -it ${MAINTAINER}/${NAME}

build: vet
	@echo Building Container
	@docker build -t ${MAINTAINER}/${NAME} .

vet:
	@echo Tidy Code
	@go mod tidy
	@echo Formatting Code
	@go fmt ./...
	@echo Vetting Code
	@go vet .

push: build
	@echo Tagging Container
	docker tag ${MAINTAINER}/${NAME}:latest ${MAINTAINER}/${NAME}:${VERSION}
	@echo Pushing Container
	docker push ${MAINTAINER}/${NAME}:latest
	@echo Pushing Container
	docker push ${MAINTAINER}/${NAME}:${VERSION}

test:
	@echo Running Unit Tests
	@mkdir -p .coverage
	@GOOS=darwin go test -tags test -coverprofile=/tmp/cov.out ./...
	@go tool cover -html=/tmp/cov.out -o=.coverage/cov.html
	@open .coverage/cov.html

tag:
	@echo Creating Git Tag
	git tag v${VERSION}
	@echo Pushing Git Tag
	git push origin --tags

NAME := logger
MAINTAINER := cachelab
VERSION := $(shell grep "const version =" main.go | cut -d\" -f2)

.PHONY: *

default: run

run: build
	docker run --network host ${MAINTAINER}/${NAME}

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
	docker tag ${MAINTAINER}/${NAME}:latest ${MAINTAINER}/${NAME}:${VERSION}
	docker push ${MAINTAINER}/${NAME}:latest
	docker push ${MAINTAINER}/${NAME}:${VERSION}

test:
	@echo Running Unit Tests
	@mkdir -p .coverage
	@GOOS=darwin go test -tags test -coverprofile=/tmp/cov.out ./...
	@go tool cover -html=/tmp/cov.out -o=.coverage/cov.html
	@open .coverage/cov.html

tag:
	git tag v${VERSION}
	git push origin --tags

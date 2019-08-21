FROM golang:1.12.9-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY logger /bin/logger

USER nobody

ENTRYPOINT ["/bin/logger"]

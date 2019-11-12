FROM golang:1.13.4-alpine

MAINTAINER Cache Lab <hello@cachelab.co>

COPY logger /bin/logger

USER nobody

ENTRYPOINT ["/bin/logger"]

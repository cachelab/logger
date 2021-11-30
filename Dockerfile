# Stage 1
FROM golang:1.17.3-alpine as builder

MAINTAINER Cache Lab <hello@cachelab.co>

ADD ./ /go/src/svc

WORKDIR /go/src/svc

RUN go build -o svc

# Stage 2
FROM alpine:3.15.0

COPY --from=builder /go/src/svc /usr/bin/

ENTRYPOINT ["/usr/bin/svc"]

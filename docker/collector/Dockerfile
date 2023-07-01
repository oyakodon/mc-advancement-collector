FROM golang:1.20.5-alpine AS build-env

ENV GO111MODULE=on

RUN apk --no-cache add git make build-base

RUN mkdir -p /build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 \
    go build -o=/build/app ./handler/collector.go

FROM alpine:3.18

RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=build-env /build/app /build/app
RUN chmod u+x /build/app

WORKDIR /build

ENTRYPOINT ["/build/app"]

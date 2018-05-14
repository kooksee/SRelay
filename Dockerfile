FROM golang:alpine as builder
MAINTAINER xtaci <daniel820313@gmail.com>
RUN apk update && \
    apk upgrade && \
    apk add git
RUN go get -ldflags "-X main.VERSION=$(date -u +%Y%m%d) -s -w" github.com/xtaci/kcptun/client && go get -ldflags "-X main.VERSION=$(date -u +%Y%m%d) -s -w" github.com/xtaci/kcptun/server

FROM alpine:3.6
COPY --from=builder /go/bin /bin
EXPOSE 29900/udp
EXPOSE 12948


FROM ubuntu:16.04

RUN rm -rf /app && mkdir /app && mkdir /kdata
COPY main /app/server
WORKDIR /app

EXPOSE 46380
EXPOSE 46381
EXPOSE 46382
EXPOSE 46383

VOLUME /kdata

CMD ["s"]
ENTRYPOINT ["/app/server","--home","/kdata"]
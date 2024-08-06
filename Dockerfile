FROM golang:1.22.4 AS build

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -o dm-db-exporter

FROM ubuntu:22.04
LABEL authors=""
LABEL maintainer=""

ENV VERSION 0.0.1
ENV DATA_SOURCE_NAME dm://SYSDBA:SYSDBA@localhost:5236?autoCommit=true

WORKDIR /opt/dm-db-exporter
COPY --from=build /build/dm-db-exporter /opt/dm-db-exporter
ADD default-metrics.toml /opt/dm-db-exporter/default-metrics.toml
RUN chmod 755 /opt/dm-db-exporter/dm-db-exporter

EXPOSE 9161

ENTRYPOINT ["/opt/dm-db-exporter/dm-db-exporter"]

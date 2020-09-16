FROM golang:1.14 AS build
WORKDIR /go/src/dmdb_exporter

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

COPY . .
RUN go build  -o  dmdb_exporter main.go


FROM frolvlad/alpine-glibc:glibc-2.29
LABEL authors=""
LABEL maintainer=""

ENV VERSION dev-0.1.0
ENV DATA_SOURCE_NAME dm://SYSDBA:SYSDBA@localhost:5236?autoCommit=true

COPY --from=build /go/src/dmdb_exporter/dmdb_exporter /dmdb_exporter
ADD ./default-metrics.toml /default-metrics.toml

RUN chmod 755 /dmdb_exporter

EXPOSE 9161

ENTRYPOINT ["/dmdb_exporter"]
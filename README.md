# DmService DB Exporter 达梦数据库Exporter


##### Table of Contents  

[Description](#description)  
[Installation](#installation)  
[Running](#running)  
[Grafana](#grafana)  
[Troubleshooting](#troubleshooting)  
[Operating principles](operating-principles.md)

# Description

A [Prometheus](https://prometheus.io/) exporter for DmService copy after the Oracle exporter. 

The following metrics are exposed currently.

- dmdb_exporter_last_scrape_duration_seconds
- dmdb_exporter_last_scrape_error
- dmdb_exporter_scrapes_total
- dmdb_session_active
- dmdb_session_max
- dmdb_session_used
- dmdb_tablespace_free_percent
- dmdb_tablespace_free_space
- dmdb_tablespace_total_space
- dmdb_up

# Installation

## Docker

You can run via Docker using an existing image. 

```bash
docker run -d --name dmdb_exporter  -p 9161:9161 -e DATA_SOURCE_NAME=dm://SYSDBA:SYSDBA@localhost:5236?autoCommit=true ${image_name}
```



## Binary Release

Pre-compiled versions for Linux 64 bit  can be found under [releases].


# Running

Ensure that the environment variable DATA_SOURCE_NAME is set correctly before starting. For Example:

```bash
# using a complete url:
export DATA_SOURCE_NAME=dm://SYSDBA:SYSDBA@localhost:5236?autoCommit=true
# Then run the exporter
/path/to/binary/dmdb_exporter --log.level error --web.listen-address 0.0.0.0:9161
```

# Integration with System D

Create file **/etc/systemd/system/dmdb_exporter.service** with the following content:

    [Unit]
    Description=Service for dm telemetry client
    After=network.target
    [Service]
    Type=oneshot
    #User=dmdb_exporter
    ExecStart=/path/of/the/dmdb_exporter --log.level error --web.listen-address 0.0.0.0:9161
    [Install]
    WantedBy=multi-user.target

Then tell System D to read files:

    systemctl daemon-reload

Start this new service:

    systemctl start dmdb_exporter

Check service status:

    systemctl status dmdb_exporter

## Usage

```bash
usage: dmdb_exporter [<flags>]

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9161"
                                 Address to listen on for web interface and telemetry. (env: LISTEN_ADDRESS)
      --web.telemetry-path="/metrics"
                                 Path under which to expose metrics. (env: TELEMETRY_PATH)
      --default.metrics="default-metrics.toml"
                                 File with default metrics in a TOML file. (env: DEFAULT_METRICS)
      --custom.metrics=""        File that may contain various custom metrics in a TOML file. (env: CUSTOM_METRICS)
      --query.timeout="5"        Query timeout (in seconds). (env: QUERY_TIMEOUT)
      --database.maxIdleConns=0  Number of maximum idle connections in the connection pool. (env: DATABASE_MAXIDLECONNS)
      --database.maxOpenConns=10
                                 Number of maximum open connections in the connection pool. (env: DATABASE_MAXOPENCONNS)
      --log.level="info"         Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                                 Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version                  Show application version.
```

# Default metrics

This exporter comes with a set of default metrics defined in **default-metrics.toml**. You can modify this file or
provide a different one using ``default.metrics`` option.

# Custom metrics

This exporter does not have the metrics you want? You can provide new one using TOML file. To specify this file to the
exporter, you can:
- Use ``--custom.metrics`` flag followed by the TOML file
- Export CUSTOM_METRICS variable environment (``export CUSTOM_METRICS=my-custom-metrics.toml``)

This file must contain the following elements:
- One or several metric section (``[[metric]]``)
- For each section a context, a request and a map between a field of your request and a comment.

Here's a simple example:

```
[[metric]]
context = "test"
request = "SELECT 1 as value_1, 2 as value_2 FROM DUAL"
metricsdesc = { value_1 = "Simple example returning always 1.", value_2 = "Same but returning always 2." }
```

This file produce the following entries in the exporter:

```
# HELP dmdb_test_value_1 Simple example returning always 1.
# TYPE dmdb_test_value_1 gauge
dmdb_test_value_1 1
# HELP dmdb_test_value_2 Same but returning always 2.
# TYPE dmdb_test_value_2 gauge
dmdb_test_value_2 2
```

You can also provide labels using labels field. Here's an example providing two metrics, with and without labels:

```
[[metric]]
context = "context_no_label"
request = "SELECT 1 as value_1, 2 as value_2 FROM DUAL"
metricsdesc = { value_1 = "Simple example returning always 1.", value_2 = "Same but returning always 2." }

[[metric]]
context = "context_with_labels"
labels = [ "label_1", "label_2" ]
request = "SELECT 1 as value_1, 2 as value_2, 'First label' as label_1, 'Second label' as label_2 FROM DUAL"
metricsdesc = { value_1 = "Simple example returning always 1.", value_2 = "Same but returning always 2." }
```

This TOML file produce the following result:

```
# HELP dmdb_context_no_label_value_1 Simple example returning always 1.
# TYPE dmdb_context_no_label_value_1 gauge
dmdb_context_no_label_value_1 1
# HELP dmdb_context_no_label_value_2 Same but returning always 2.
# TYPE dmdb_context_no_label_value_2 gauge
dmdb_context_no_label_value_2 2
# HELP dmdb_context_with_labels_value_1 Simple example returning always 1.
# TYPE dmdb_context_with_labels_value_1 gauge
dmdb_context_with_labels_value_1{label_1="First label",label_2="Second label"} 1
# HELP dmdb_context_with_labels_value_2 Same but returning always 2.
# TYPE dmdb_context_with_labels_value_2 gauge
dmdb_context_with_labels_value_2{label_1="First label",label_2="Second label"} 2
```

Last, you can set metric type using **metricstype** field.

```
[[metric]]
context = "context_with_labels"
labels = [ "label_1", "label_2" ]
request = "SELECT 1 as value_1, 2 as value_2, 'First label' as label_1, 'Second label' as label_2 FROM DUAL"
metricsdesc = { value_1 = "Simple example returning always 1 as counter.", value_2 = "Same but returning always 2 as gauge." }
# Can be counter or gauge (default)
metricstype = { value_1 = "counter" }
```

This TOML file will produce the following result:

```
# HELP dmdb_test_value_1 Simple test example returning always 1 as counter.
# TYPE dmdb_test_value_1 counter
dmdb_test_value_1 1
# HELP dmdb_test_value_2 Same test but returning always 2 as gauge.
# TYPE dmdb_test_value_2 gauge
dmdb_test_value_2 2
```

# Customize metrics in a docker image

If you run the exporter as a docker image and want to customize the metrics, you can use the following example:

```Dockerfile
FROM xxx/dmdb_exporter:latest

COPY custom-metrics.toml /

ENTRYPOINT ["/dmdb_exporter", "--custom.metrics", "/custom-metrics.toml"]
```

# Integration with Grafana

An example Grafana dashboard is available [here](https://grafana.com/dashboards/3333).

# Build

## Docker build

To build  Alpine image, run the following command:

    make docker


## Linux binaries

Retrieve Linux binaries:

    go build main.go 


# FAQ/Troubleshooting

## Unable to convert current value to float (metric=par,metri...in.go:285

DmService is trying to send a value that we cannot convert to float. This could be anything like 'UNLIMITED' or 'UNDEFINED' or 'WHATEVER'.

In this case, you must handle this problem by testing it in the SQL request. Here an example available in default metrics:

```toml
[[metric]]
context = "resource"
labels = [ "resource_name" ]
metricsdesc = { current_utilization= "Generic counter metric from v$resource_limit view in Oracle (current value).", limit_value="Generic counter metric from v$resource_limit view in Oracle (UNLIMITED: -1)." }
request="SELECT resource_name,current_utilization,CASE WHEN TRIM(limit_value) LIKE 'UNLIMITED' THEN '-1' ELSE TRIM(limit_value) END as limit_value FROM v$resource_limit"
```

If the value of limite_value is 'UNLIMITED', the request send back the value -1.

You can increase the log level (`--log.level debug`) in order to get the statement generating this error.

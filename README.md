# YugabyteDB client for go

Work in progress. Current state: this can definitely work.

## Start YugabyteDB in a container

```sh
docker run \
    --rm \
    --hostname localhost \
    -ti \
        -p 7000:7000 \
        -p 7100:7100 \
        -p 5433:5433 \
        local/yugabytedb:2.7.2.0 \
        /home/yugabyte/bin/yb-master \
            --master_addresses=localhost:7100 \
            --rpc_bind_addresses=0.0.0.0:7100 \
            --fs_data_dirs=/tmp \
            --logtostderr --replication_factor=1
```

## Run the client

```
go run ./main.go [command] [flags]
```

where the command is one of:

- `get-master-registration`
- `is-load-balanced`
- `list-masters`
- `list-tablet-servers`

### Flags

Common flags:

- `--master`: host port of the master to query, default `127.0.0.1:7100`
- `--operation-timeout`: RPC operation timeout, duration string (`5s`, `1m`, ...), default `5s`
- `--tls-ca-cert-file-path`: full path to the CA certificate file, default `empty string`
- `--tls-cert-file-path`: full path to the certificate file, default `empty string`
- `--tls-key-file-path`: full path to the key file, default `empty string`

Logging flags:

- `--log-level`: log level, default `info`
- `--log-as-json`: log entries as JSON, default `false`
- `--log-color`: log colored output, default `false`
- `--log-force-color`: force colored output, default `false`

#### Command specific flags

##### is-load-balanced

- `--expected-num-servers`: int32, how many servers to include in this check, default `-1` (`undefined`)

##### list-tablet-servers

- `--primary-only`: boolean, list primary tablet servers only, default `false`

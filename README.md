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

- `check-exists`: Check that a table exists.
- `get-load-move-completion`: Get the completion percentage of tablet load move from blacklisted servers.
- `get-master-registration`: Get master registration info.
- `get-universe-config`: Get the placement info and blacklist info of the universe.
- `is-load-balanced`: Check if master leader thinks that the load is balanced across TServers.
- `is-server-ready`: Check if server is ready to serve IO requests.
- `list-masters`: List all the masters in this database.
- `list-tablet-servers`: List all the tablet servers in this database.
- `ping`: Ping a certain YB server.

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

##### check-exists

- `--keyspace`: string, keyspace name to check in, default `<empty string>`
- `--name`: string, table name to check for, default `<empty string>`
- `--uuid`: string, table identified (uuid) to check for, default `<empty string>`

##### is-load-balanced

- `--expected-num-servers`: int32, how many servers to include in this check, default `-1` (`undefined`)

##### list-tablet-servers

- `--primary-only`: boolean, list primary tablet servers only, default `false`

##### ping

- `--host`: string, host to ping, default `<empty string>`
- `--port`: int, port to ping, default `0`, must be higher than `0`

##### is-server-ready

- `--host`: string, host to check, default `<empty string>`
- `--port`: int, port to check, default `0`, must be higher than `0`
- `--is-tserver`: boolean, when `true` - indicated a TServer, default `false`

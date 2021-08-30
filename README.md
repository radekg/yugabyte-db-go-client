# YugabyteDB client for go

Work in progress. Current state: this can definitely work.

## Go client

### Usage

The `github.com/radekg/yugabyte-db-go-client/client/cli` provides a reference client implementation.

**TL;DR**: here's how to use the API client directly from _go_:

```go
package main

import (
    "github.com/radekg/yugabyte-db-go-client/client/base"
    "github.com/radekg/yugabyte-db-go-client/configs"
    ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
    
    "github.com/hashicorp/go-hclog"

    "encoding/json"
    "fmt"
    "time"
)

func main() {

    // construct the configuration:
    cfg := &configs.YBClientConfig{
        MasterHostPort: "127.0.0.1:7100",
        OpTimeout:      uint32(5000),
        // use TLSConfig of type *tls.Config to configure TLS
    }

    // create a logger:
    logger := hclog.Default()

    // create a client:
    connectedClient, err := base.Connect(cfg, logger)

    if err != nil {
        panic(err)
    }

    // wait for connection status:
    select {
    case err := <-connectedClient.OnConnectError():
        logger.Error("failed connecting a client", "reason", err)
        panic(err)
    case <-connectedClient.OnConnected():
        logger.Debug("client connected")
    }
    // when not using panics further down, use
    // defer connectedClient.Close()

    // create the request payload:
    payload := &ybApi.ListMastersRequestPB{}
    // create the response payload, it will be populated with the response data
    // if the request succeeded:
    responsePayload := &ybApi.ListMastersResponsePB{}

    // execute the request:
    if err := connectedClient.Execute(payload, responsePayload); err != nil {
        connectedClient.Close()
        logger.Error("failed executing the request", "reason", err)
        panic(err)
    }

    // some of the payloads provide their own error responses,
    // handle it like this:
    if err := responsePayload.GetError(); err != nil {
        connectedClient.Close()
        logger.Error("request returned an error", "reason", err)
        panic(err)
    }

    // do something with the result:
    bytes, err := json.MarshalIndent(responsePayload, "", "  ")
    if err != nil {
        connectedClient.Close()
        logger.Error("failed marshalling the response as JSON", "reason", err)
        panic(err)
    }

    fmt.Println(string(bytes))

    // close the client at the very end
    connectedClient.Close()

}
```

## CLI client

### Start YugabyteDB in a container

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

The goal is to provide a client implementing the functionality of the official YugabyteDB Java client.

### Usage

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
- `leader-step-down`: Try to force the current leader to step down, requires `--destination-uuid` and `--tablet-id`.
- `list-masters`: List all the masters in this database.
- `list-tablet-servers`: List all the tablet servers in this database.
- `master-leader-step-down`: Try to force the current master leader to step down.
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

### Command specific flags

#### check-exists

- `--keyspace`: string, keyspace name to check in, default `<empty string>`
- `--name`: string, table name to check for, default `<empty string>`
- `--uuid`: string, table identified (uuid) to check for, default `<empty string>`

#### leader-step-down

- `--destination-uuid`: UUID of server this request is addressed to, default `empty` - not specified
- `--disable-graceful-transition`: `boolean`, if `new-leader-uuid` is not specified, the current leader will attempt to gracefully transfer leadership to another peer; setting this flag disables that behavior, default `false`
- `--new-leader-uuid`: UUID of the server that should run the election to become the new leader, default `empty` - not specified
- `--tablet-id`: the id of the tablet, default `empty` - not specified

#### list-tablet-servers

- `--primary-only`: boolean, list primary tablet servers only, default `false`

#### ping

- `--host`: string, host to ping, default `<empty string>`
- `--port`: int, port to ping, default `0`, must be higher than `0`

#### is-load-balanced

- `--expected-num-servers`: int32, how many servers to include in this check, default `-1` (`undefined`)

#### is-server-ready

- `--host`: string, host to check, default `<empty string>`
- `--port`: int, port to check, default `0`, must be higher than `0`
- `--is-tserver`: boolean, when `true` - indicated a TServer, default `false`

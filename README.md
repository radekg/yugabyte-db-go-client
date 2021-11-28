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

### Usage

```
go run ./main.go [command] [flags]
```

where the command is one of:

- `check-exists`: Check that a table exists.
- `describe-table`: Info on a table in this database.
- `get-load-move-completion`: Get the completion percentage of tablet load move from blacklisted servers.
- `get-master-registration`: Get master registration info.
- `get-tablets-for-table`: Fetch tablet information for a given table.
- `get-universe-config`: Get the placement info and blacklist info of the universe.
- `is-load-balanced`: Check if master leader thinks that the load is balanced across TServers.
- `is-server-ready`: Check if server is ready to serve IO requests.
- `leader-step-down`: Try to force the current leader to step down, requires `--destination-uuid` and `--tablet-id`.
- `list-masters`: List all the masters in this database.
- `list-tables`: List all tables in this database.
- `list-tablet-servers`: List all the tablet servers in this database.
- `master-leader-step-down`: Try to force the current master leader to step down.
- `ping`: Ping a certain YB server.
- `set-load-balancer-state`: Set the load balancer state.

#### Snapshot commands

- `create-snapshot`: Creates a snapshot of an entire keyspace or selected tables in a keyspace.
- `list-snapshots`: List snapshots.

### Flags

Common flags:

- `--master`: host port of the master to query, default `127.0.0.1:7100`
- `--operation-timeout`: RPC operation timeout, duration string (`5s`, `1m`, ...), default `60s`
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

#### describe-table

- `--keyspace`: string, keyspace name to check in, default `<empty string>`, ignored when using `--uuid`
- `--name`: string, table name to check for, default `<empty string>`
- `--uuid`: string, table identified (uuid) to check for, default `<empty string>`

Examples:

- describe table `test` in the `yugabyte` database: `cli describe-table --keyspace yugabyte --name test`
- describe table with ID `000033c0000030008000000000004000`: `cli describe-table --uuid 000033c0000030008000000000004000`

#### get-tablets-for-table

- `--keyspace`: string, keyspace to describe the table in, default `empty string`
- `--name`: string, table name to check for, default `empty string`
- `--uuid`: string, table identifier to check for, default `empty string`
- `--partition-key-start`: base64 encoded, partition key range start, default `empty`
- `--partition-key-end`: base64 encoded, partition key range end, default `empty`
- `--max-returned-locations`: uint32, maximum number of returned locations, default `10`
- `--require-tablet-running`: boolean, require tablet running, default `false`

#### leader-step-down

- `--destination-uuid`: UUID of server this request is addressed to, default `empty` - not specified
- `--disable-graceful-transition`: boolean, if `new-leader-uuid` is not specified, the current leader will attempt to gracefully transfer leadership to another peer; setting this flag disables that behavior, default `false`
- `--new-leader-uuid`: UUID of the server that should run the election to become the new leader, default `empty` - not specified
- `--tablet-id`: the id of the tablet, default `empty` - not specified

#### list-tables

- `--name-filter`: string, When used, only returns tables that satisfy a substring match on `name_filter`, default `empty string`
- `--keyspace`: string, the namespace name to fetch info, default `empty string`
- `--exclude-system-tables`: boolean, exclude system tables, default `false`
- `--include-not-running`: boolean, include not running, default `false`
- `--relation-type`: list of strings, filter tables based on RelationType - supported values: `system_table`, `user_table`, `index`, default: all values

Examples:

- list all PostgreSQL `system_platform` relations: `cli list-tables --keyspace ysql.system_platform`
- list all PostgreSQL `postgres` relations: `cli list-tables --keyspace ysql.postgres`
- list all PostgreSQL `yugabyte` relations: `cli list-tables --keyspace ysql.yugabyte`
- list all PostgreSQL `template0` relations: `cli list-tables --keyspace ysql.template0`
- list all CQL `system_schema` relations: `cli list-tables --keyspace ycql.system_schema`
- list all Redis `system_redis` relations: `cli list-tables --keyspace yedis.system_redis`

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

#### set-load-balancer-state

Options are mutually exclusive, exactly one has to be set:

- `--enabled`: boolean, default `false`, new desired state: enabled
- `--disabled`: boolean, default `false`, new desired state: disabled

#### Snapshot commands

##### create-snapshot

Multiple `--name` and `--uuid` values can be combined together. To create a snapshot of an entire keyspace, do not specify any `--name` or `--uuid`.

- `--keyspace`: string, keyspace name to create snapshot of, default `<empty string>`
- `--name`: repeated string, table name to create snapshot of, default `empty list`
- `--uuid`: repeated string, table ID to create snapshot of, default `empty list`
- `--transaction-aware`: boolean, transaction aware, default `false`
- `--add-indexes`: boolean, add indexes, default `false`, YCQL only
- `--imported`: boolean, interpret this snapshot as imported, default `false`
- `--schedule-id`: base64 encoded, create snapshot to this schedule, other fields are ignored, default `empty`

Examples:

- create a snapshot of an entire YSQL `yugabyte` database: `cli create-snapshot --keyspace ysql.yugabyte`
- create a snapshot of selected YSQL tables in the `yugabyte` database: `cli create-snapshot --keyspace ysql.yugabyte --name table --name another-table`
- create a snapshot of selected YCQL tables in the `example` database: `cli create-snapshot --keyspace ycql.example --name table --add-indexes`

##### list-snapshots

- `--snapshot-id`: string, Snapshot identifier, default `empty string` (not defined)
- `--list-deleted-snapshots`: boolean, list deleted snapshots, default `false`
- `--prepare-for-backup`: boolean, prepare for backup, default `false`

## Minimal YugabyteDB cluster in Docker compose

This repository contains a minimal YugabyteDB Docker compose setup which can be used for client testing or validation.

To start the cluster:

```sh
cd .compose/
docker compose -f yugabytedb-minimal.yml up
```

To restart:

```sh
docker compose -f yugabytedb-minimal.yml rm
docker compose -f yugabytedb-minimal.yml up
```

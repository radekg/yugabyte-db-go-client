# YugabyteDB client for go

Work in progress. Current state: this can definitely work.

## Go client

### Usage

The `github.com/radekg/yugabyte-db-go-client/client/implementation` provides a reference client implementation.

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
- `delete-snapshot`: Delete a snapshot.
- `export-snapshot`: Exports a snapshot.
- `import-snapshot`: Imports a snapshot.
- `list-snapshots`: List snapshots.
- `list-snapshot-restorations`: List snapshot restorations.
- `restore-snapshot`: Restore a snapshot.
- `restore-snapshot-schedule`: Restore a snapshot schedule.

- `create-snapshot-schedule`: Creates a snapshot schedule from an entire keyspace or selected tables in the keyspace.
- `delete-snapshot-schedule`: Delete a snapshot schedule.
- `list-snapshot-schedules`: List snapshot schedules.

### Flags

Common flags:

- `--master`: string, repeated, host port of the master to query, default `127.0.0.1:7100, 127.0.0.1:7101, 127.0.0.1:7102`
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
- `--relation-type`: list of strings, filter tables based on RelationType - supported values: `system_table`, `user_table`, `index_table`, default: all values

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

- `--keyspace`: string, keyspace name to create snapshot of, default `<empty string>`
- `--name`: repeated string, table name to create snapshot of, default `empty list`
- `--uuid`: repeated string, table ID to create snapshot of, default `empty list`
- `--schedule-id`: base64 encoded, create snapshot to this schedule, other fields are ignored, default `empty`
- `--base64-encoded`: boolean, base64 decode given schedule ID before handling over to the API, default `false`

Remarks:

- Multiple `--name` and `--uuid` values can be combined together.
- YSQL keyspace snapshots do not support explicit `--name` and `--uuid` selection.
- To create a snapshot of an entire keyspace, do not specify any `--name` or `--uuid`. YCQL only.
- `yedis.*` keyspaces are not supported.

Examples:

- create a snapshot of an entire YSQL `yugabyte` database: `cli create-snapshot --keyspace ysql.yugabyte`
- create a snapshot of selected YCQL tables in the `example` database: `cli create-snapshot --keyspace ycql.example --name table`

##### delete-snapshot

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)

##### list-snapshots

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)
- `--list-deleted-snapshots`: boolean, list deleted snapshots, default `false`
- `--prepare-for-backup`: boolean, prepare for backup, default `false`

##### list-snapshot-restorations

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)
- `--restoration-id`: string, restoration identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)

##### export-snapshot

- `--snapshot-id`: string, snapshot identifier- literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--file-path`: string, full path to the export file, parent directories must exist, default `empty`

##### import-snapshot

- `--file-path`: string, full path to the exported snapshot file
- `--keyspace`: string, fully qualified keyspace name, for example `ycql.system_namespace`, no effect for YSQL imports, default `empty`
- `--table-name`: string, repeated, table name to import, no effect for YSQL snapshots, default `empty list`

##### restore-snapshot

- `--schedule-id`: string, schedule identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--restore-target`: exact past HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), absolute Timing Option: Max HybridTime, or relative past interval, default `empty` (undefined)

##### restore-snapshot-schedule

- `--snapshot-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)
- `--restore-target`: exact past HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), absolute Timing Option: Max HybridTime, or relative past interval, default `empty` (undefined)

##### create-snapshot-schedule

- `--keyspace`: string, keyspace name to create snapshot of, default `<empty string>`
- `--interval`: duration expression (`1h`, `1d`, ...), interval for taking snapshot in seconds, default `0` (undefined)
- `--retention-duration`: duration expression (`1h`, `1d`, ...), how long store snapshots in seconds, default `0` (undefined)
- `--delete-after`: exact future HT (`16 digit literal`) or duration expression (`1h`, `5h15m`, ...), how long until schedule is removed in seconds, hybrid time will be calculated by fetching server hybrid time and adding this value, default `0` (undefined)

Examples:

- create a snapshot schedule of an entire YSQL `yugabyte` database: `cli create-snapshot-schedule --keyspace ysql.yugabyte --interval 1h --retention-duration 2h --delete-after 1h`
- create a snapshot schedule of selected YSQL tables in the `yugabyte` database: `cli create-snapshot-schedule --keyspace ysql.yugabyte --name table --name another-table`

##### delete-snapshot-schedule

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, required, default `empty string` (not defined)

##### list-snapshot-schedules

- `--schedule-id`: string, snapshot identifier - literal ID or Base64 encoded value from YugabyteDB RPC API, optional, default `empty string` (not defined)

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

## Docker image

Build the Docker image:

```sh
make docker-image
```

Run against the provided minimal YugabyteDB cluster:

```sh
docker run --rm \
    --net yb-client-minimal \
    -ti local/ybdb-go-cli:0.0.1 \
    list-masters --master yb-master-1:7100 \
                 --master yb-master-2:7100 \
                 --master yb-master-3:7100
```

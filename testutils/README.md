# YugabyteDB Test Kit

YugabyteDB for embedding in go tests. Depends on Docker.

## Usage

Individual packages contain tests showing exact usage patterns for your own tests.

## Example

Run three masters with three TServers inside of the test and query YSQL on one of the TServers:

```go
package myprogram

import (
    "testing"
    "time"

    "github.com/hashicorp/go-hclog"
    "github.com/radekg/yugabyte-db-go-client/testutils/common"
    "github.com/radekg/yugabyte-db-go-client/testutils/master"
    "github.com/radekg/yugabyte-db-go-client/testutils/tserver"
    "github.com/radekg/yugabyte-db-go-client/client/implementation"
    "github.com/radekg/yugabyte-db-go-client/configs"
)

func TestClusterIntegration(t *testing.T) {

    masterTestCtx := master.SetupMasters(t, &common.TestMasterConfiguration{
        ReplicationFactor: 3,
    })
    defer masterTestCtx.Cleanup()

    client, err := implementation.MasterLeaderConnectedClient(&configs.CliConfig{
        MasterHostPort: masterTestCtx.MasterExternalAddresses(),
        OpTimeout:      time.Duration(time.Second * 5),
    }, hclog.Default())

    if err != nil {
        t.Fatal(err)
    }
    defer client.Close()

    listMastersPb, err := client.ListMasters()
    if err != nil {
        t.Fatal(err)
    }

    t.Log(" ========> ", listMastersPb)

    // start a TServer:
    tserver1Ctx := tserver.SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
        TServerID: "my-tserver-1",
    })
    defer tserver1Ctx.Cleanup()

    // start a TServer:
    tserver2Ctx := tserver.SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
        TServerID: "my-tserver-2",
    })
    defer tserver2Ctx.Cleanup()

    // start a TServer:
    tserver3Ctx := tserver.SetupTServer(t, masterTestCtx, &common.TestTServerConfiguration{
        TServerID: "my-tserver-3",
    })
    defer tserver3Ctx.Cleanup()

    listTServersPb, err := client.ListTabletServers(&configs.OpListTabletServersConfig{})
    if err != nil {
        t.Fatal(err)
    }

    t.Log(" ==> Received TServer list", listTServersPb)

}
```

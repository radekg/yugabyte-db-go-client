# YugabyteDB client for go

Work in progress. Current state: this can definitely work.

## Go client

### Usage

The `github.com/radekg/yugabyte-db-go-client/client/implementation` provides a reference client implementation.

**TL;DR**: here's how to use the API client directly from _go_:

```go
package main

import (
    "github.com/radekg/yugabyte-db-go-client/client"
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

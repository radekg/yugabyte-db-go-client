# YugabyteDB client for go

Work in progress. Current state: it does what it says it does.

## Go client

### Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/hashicorp/go-hclog"

    "github.com/radekg/yugabyte-db-go-client/client"
    "github.com/radekg/yugabyte-db-go-client/configs"
    "github.com/radekg/yugabyte-db-go-client/errors"

    ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"
)

func main() {

    // construct the configuration:
    cfg := &configs.YBClientConfig{
        MasterHostPort: []string{"127.0.0.1:7100", "127.0.0.1:17000", "127.0.0.1:27000"},
        OpTimeout:      time.Duration(time.Second * 5),
    }

    customLogger := hclog.Default()

    client := client.NewYBClient(cfg).
        WithLogger(customLogger.Named("custom-client-logger"))

    if err := client.Connect(); err != nil {
        panic(err)
	}

    request := &ybApi.ListMastersRequestPB{}
    response := &ybApi.ListMastersResponsePB{}
    err := client.Execute(request, response)
    if err != nil {
        client.Close()
        logger.Error("failed executing the request", "reason", err)
        panic(err)
    }

    // some of the payloads provide their own error responses,
    // handle it like this:
    if err := response.GetError(); err != nil {
        client.Close()
        logger.Error("request returned an error", "reason", err)
        panic(err)
    }

    // do something with the result:
    bytes, err := json.MarshalIndent(response, "", "  ")
    if err != nil {
        client.Close()
        logger.Error("failed marshalling the response as JSON", "reason", err)
        panic(err)
    }

    fmt.Println("successful masters response", string(bytes))

    // close the client at the very end
    client.Close()

}
```

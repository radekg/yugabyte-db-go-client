package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"

	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"

	"github.com/radekg/yugabyte-db-go-client/client"
	"github.com/radekg/yugabyte-db-go-client/utils"
)

func getMasterRegistration(conn net.Conn) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(0),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("GetMasterRegistration"),
		},
		TimeoutMillis: utils.PUint32(5000),
	}
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	b := bytes.NewBuffer([]byte{})
	utils.WriteMessages(b, requestHeader, payload)
	conn.Write(b.Bytes())
}

func main() {

	queryMasters := false
	queryMasterRegistration := false

	cfg := &client.YBClientConfig{
		MasterHostPort: "127.0.0.1:7100",
	}

	connectedClient, err := client.Connect(cfg)
	if err != nil {
		panic(err)
	}
	select {
	case err := <-connectedClient.OnConnectError():
		fmt.Println("client failed to connect, reason", err)
		os.Exit(1)
	case <-connectedClient.OnConnected():
	}

	defer connectedClient.Close()
	fmt.Println(" ====> client is now connected")

	// list masters:
	if queryMasters {
		masters, err := connectedClient.ListMasters()
		if err != nil {
			fmt.Println(" ====> failed reading masters", err)
			return
		}
		jsonBytes, err := json.MarshalIndent(masters, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes))
	}

	// get master registration:
	if queryMasterRegistration {
		registration, err := connectedClient.GetMasterRegistration()
		if err != nil {
			fmt.Println(" ====> failed reading master registration", err)
			return
		}
		jsonBytes2, err := json.MarshalIndent(registration, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes2))
	}

	// list tablet servers:
	tablets, err := connectedClient.ListTabletServers()
	if err != nil {
		fmt.Println(" ====> failed reading tablet servers", err)
		return
	}
	if len(tablets.Servers) == 0 {
		fmt.Println(" ====> no tablet servers present")
	} else {
		jsonBytes3, err := json.MarshalIndent(tablets, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonBytes3))
	}

	/*
		conn, err := net.Dial("tcp")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// Once connected, send the YugabyteDB header:
		// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/TabletClient.java#L593
		header := append([]byte("YB"), 1)
		conn.Write(header)
		// this doesn't reply with anything:
		mustReadFromConn(conn)

		//
		getMasterRegistration(conn)
		mustReadFromConn(conn)

		listMasters(conn)
		readListMastersResponse(mustReadFromConn(conn))
	*/
}

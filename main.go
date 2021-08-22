package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
	ybApi "github.com/radekg/yugabyte-db-go-proto/v2/yb/api"

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

func listMasters(conn net.Conn) {
	requestHeader := &ybApi.RequestHeader{
		CallId: utils.PInt32(1),
		RemoteMethod: &ybApi.RemoteMethodPB{
			ServiceName: utils.PString("yb.master.MasterService"),
			MethodName:  utils.PString("ListMasters"),
		},
		TimeoutMillis: utils.PUint32(5000),
	}
	payload := &ybApi.GetMasterRegistrationRequestPB{}
	b := bytes.NewBuffer([]byte{})
	utils.WriteMessages(b, requestHeader, payload)
	conn.Write(b.Bytes())
}

func mustReadFromConn(conn net.Conn) []byte {
	buf2 := make([]byte, 1024*1024)
	read2, err2 := conn.Read(buf2)
	if err2 != nil {
		panic(err2)
	}
	return buf2[0:read2]
}

func readListMastersResponse(input []byte) {

	reader := bytes.NewReader(input)

	// Read the complete data length:
	dataLength, err := utils.ReadInt(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("DEBUG: the response data length is: ", dataLength)

	responseHeaderLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		panic(err)
	}

	// Now I can read the response header:
	responseHeaderBuf := make([]byte, responseHeaderLength)
	n, err := reader.Read(responseHeaderBuf)
	if err != nil {
		panic(err)
	}
	if uint64(n) != responseHeaderLength {
		panic(fmt.Errorf("expected to read %d but read %d", responseHeaderLength, n))
	}

	responseHeader := &ybApi.ResponseHeader{}
	protoErr := proto.Unmarshal(responseHeaderBuf, responseHeader)
	if protoErr != nil {
		panic(protoErr)
	}

	fmt.Println(fmt.Sprintf("DEBUG: Response to call id: %d, is error: %v, # of sidecars: %d",
		*responseHeader.CallId,
		*responseHeader.IsError,
		len(responseHeader.SidecarOffsets)))

	responsePayloadLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		panic(err)
	}

	responsePayloadBuf := make([]byte, responsePayloadLength)
	n, err = reader.Read(responsePayloadBuf)
	if err != nil {
		panic(err)
	}
	if uint64(n) != responsePayloadLength {
		panic(fmt.Errorf("expected to read %d but read %d", responsePayloadLength, n))
	}

	responsePayload := &ybApi.ListMastersResponsePB{}
	protoErr2 := proto.Unmarshal(responsePayloadBuf, responsePayload)
	if protoErr2 != nil {
		panic(protoErr2)
	}

	jsonBytes, err := json.MarshalIndent(responsePayload.Masters, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:7100")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	header := append([]byte("YB"), 1)
	conn.Write(header)
	mustReadFromConn(conn)

	getMasterRegistration(conn)
	mustReadFromConn(conn)

	listMasters(conn)
	readListMastersResponse(mustReadFromConn(conn))
}

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
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L71
	dataLength, err := utils.ReadInt(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println("DEBUG: the response data length is: ", dataLength)

	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L76
	responseHeaderLength, _, err := utils.ReadVarint(reader)
	if err != nil {
		panic(err)
	}

	// Now I can read the response header:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L78
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

	// This here is currently a guess but I believe the corretc mechanism sits here:
	// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/CallResponse.java#L113
	// The encoding/binary.ReadUvarint and encoding/binary.ReadVarint doesn't do what it supposed to do
	// hence the custom code here.
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

	// Once connected, send the YugabyteDB header:
	// https://github.com/yugabyte/yugabyte-db/blob/master/java/yb-client/src/main/java/org/yb/client/TabletClient.java#L593
	header := append([]byte("YB"), 1)
	conn.Write(header)
	// this doesn't reply with anything:
	mustReadFromConn(conn)

	//
	getMasterRegistration(conn)
	mustReadFromConn(conn)

	listMasters(conn)
	readListMastersResponse(mustReadFromConn(conn))
}

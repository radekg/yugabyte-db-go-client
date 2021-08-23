package utils

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ReadVarintMaxBufferSize is the maximum varint buffer size.
const ReadVarintMaxBufferSize = 9

// ReadVarintByteOrder byte order for varint reader.
var ReadVarintByteOrder = binary.LittleEndian

// ReadVarint reads a varint from a reader.
func ReadVarint(r io.Reader) (uint64, int, error) {
	var bufarray [ReadVarintMaxBufferSize]byte
	buf := bufarray[:]
	var value uint64
	i := 0
	i, err := io.ReadFull(r, buf[0:1])
	if err != nil {
		return 0, i, err
	}
	switch buf[0] {
	default:
		value = uint64(buf[0])
		i = 1
	case 0xfd:
		_, err := io.ReadFull(r, buf[0:2])
		if err != nil {
			return 0, i, err
		}

		value = uint64(ReadVarintByteOrder.Uint16(buf))
		i = 3
	case 0xfe:
		_, err := io.ReadFull(r, buf[0:4])
		if err != nil {
			return 0, i, err
		}
		value = uint64(ReadVarintByteOrder.Uint32(buf))
		i = 5
	case 0xff:
		_, err := io.ReadFull(r, buf[0:8])
		if err != nil {
			return 0, i, err
		}
		value = ReadVarintByteOrder.Uint64(buf)
		i = 9
	}
	return value, i, nil
}

// WriteMessages writes a variable number of protobuf messages into a given writer.
// This code is essentially based on:
// https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L274
// In the Java code, the Message header is the requestHeader,
// the Message pb is the actual payload.
func WriteMessages(b io.Writer, msgs ...protoreflect.ProtoMessage) error {
	serialized := [][]byte{}
	totalSize := 0
	// calculate the total size:
	for _, m := range msgs {
		bys, err := proto.Marshal(m)
		if err != nil {
			return err
		}
		bysLength := len(bys)
		totalSize = totalSize + bysLength
		totalSize = totalSize + computeUInt32SizeNoTag(bysLength)
		serialized = append(serialized, bys)
	}
	// write the total size to the buffer
	// per: https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L279
	if err := writeInt(b, totalSize); err != nil {
		return err
	}
	for _, s := range serialized {
		// calculate and write the individual message length varint, per:
		// - https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L282
		// - https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L285
		varintBytes := toVarintByte(len(s))
		n, err := b.Write(varintBytes)
		if err != nil {
			return err
		}
		if n != len(varintBytes) {
			return io.EOF
		}
		// write the actual payload, per:
		// - https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L283
		// - https://github.com/yugabyte/yugabyte-db/blob/v2.7.2/java/yb-client/src/main/java/org/yb/client/YRpc.java#L286
		n2, err := b.Write(s)
		if err != nil {
			return err
		}
		if n2 != len(s) {
			return io.EOF
		}
	}
	return nil
}

// PInt32 returns a pointer to an int32.
func PInt32(a int32) *int32 {
	return &a
}

// PUint32 returns a pointer to an uint32.
func PUint32(a uint32) *uint32 {
	return &a
}

// PString returns a pointer to a string.
func PString(a string) *string {
	return &a
}

// ReadInt reads an int from a reader.
func ReadInt(reader io.Reader) (int, error) {
	intBuf := make([]byte, 4)
	n, err := reader.Read(intBuf)
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, io.EOF
	}
	return int(binary.BigEndian.Uint32(intBuf)), nil
}

func writeInt(w io.Writer, i int) error {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(i))
	n, err := w.Write(arr)
	if err != nil {
		return err
	}
	if n != 4 {
		return io.EOF
	}
	return nil
}

func toVarintByte(i int) []byte {
	arr := make([]byte, 32)
	n := binary.PutUvarint(arr, uint64(i))
	return arr[0:n]
}

// This is verbatim copy of:
// https://github.com/protocolbuffers/protobuf/blob/v3.5.1/java/core/src/main/java/com/google/protobuf/CodedOutputStream.java#L728
// YugabyteDB Java client 2.7.2 depends on protobuf-java 3.5.1.
func computeUInt32SizeNoTag(value int) int {
	if (value & (^0 << 7)) == 0 {
		return 1
	}
	if (value & (^0 << 14)) == 0 {
		return 2
	}
	if (value & (^0 << 21)) == 0 {
		return 3
	}
	if (value & (^0 << 28)) == 0 {
		return 4
	}
	return 5
}

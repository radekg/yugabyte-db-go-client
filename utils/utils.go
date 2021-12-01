package utils

import (
	"encoding/binary"
	"errors"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var errOverflow32 = errors.New("binary: varint overflows a 64-bit integer")

// ReadUvarint32 reads an encoded unsigned 32-bit integer from r and returns it as a uint64.
func ReadUvarint32(r io.ByteReader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; i < binary.MaxVarintLen32; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i == binary.MaxVarintLen32-1 && b > 1 {
				return x, errOverflow32
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return x, errOverflow32
}

// ReadVarint32 reads an encoded signed integer from r and returns it as an int64.
func ReadVarint32(r io.ByteReader) (int64, error) {
	ux, err := ReadUvarint32(r) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
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

// PBool returns a pointer to an bool.
func PBool(a bool) *bool {
	return &a
}

// PInt32 returns a pointer to an int32.
func PInt32(a int32) *int32 {
	return &a
}

// PUint32 returns a pointer to an uint32.
func PUint32(a uint32) *uint32 {
	return &a
}

// PUint64 returns a pointer to an uint64.
func PUint64(a uint64) *uint64 {
	return &a
}

// PString returns a pointer to a string.
func PString(a string) *string {
	if a == "" {
		return nil
	}
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

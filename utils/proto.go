package utils

import "google.golang.org/protobuf/proto"

// DeserializeProto deserializes proto message from bytes.
func DeserializeProto(bys []byte, msg proto.Message) error {
	return proto.Unmarshal(bys, msg)
}

// SerializeProto serializes proto message to bytes.
func SerializeProto(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

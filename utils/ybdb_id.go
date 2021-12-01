package utils

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// DecodeAsYugabyteID returns a snapshot ID or an error.
func DecodeAsYugabyteID(input string, base64Encoded bool) (string, error) {
	decodedID := input
	if base64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(decodedID)
		if err != nil {
			return "", err
		}
		protoReadySnapshotID, err := ProtoYugabyteIDToString(decoded)
		if err != nil {
			return "", err
		}
		decodedID = protoReadySnapshotID
	}
	return decodedID, nil
}

// ProtoYugabyteIDToString converts the spanshot id represented as bytes to a string UUID.
func ProtoYugabyteIDToString(input []byte) (string, error) {
	aUUID := uuid.New()
	if err := aUUID.UnmarshalBinary(input); err != nil {
		return "", err
	}
	return aUUID.String(), nil
}

// StringUUIDToProtoYugabyteID converts a string UUID to the snapshot ID bytes for protobuf operations.
func StringUUIDToProtoYugabyteID(input string) ([]byte, error) {
	aUUID, err := uuid.Parse(input)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := aUUID.MarshalBinary()
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

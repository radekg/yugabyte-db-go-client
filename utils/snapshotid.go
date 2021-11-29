package utils

import (
	"encoding/base64"

	"github.com/google/uuid"
)

// SnapshotID returns a snapshot ID or an error.
func SnapshotID(input string, base64Encoded bool) (string, error) {
	givenSnapshotID := input
	if base64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(givenSnapshotID)
		if err != nil {
			return "", err
		}
		protoReadySnapshotID, err := ProtoSnapshotIDToString(decoded)
		if err != nil {
			return "", err
		}
		givenSnapshotID = protoReadySnapshotID
	}
	return givenSnapshotID, nil
}

// ProtoSnapshotIDToString converts the spanshot id represented as bytes to a string UUID.
func ProtoSnapshotIDToString(input []byte) (string, error) {
	aUUID := uuid.New()
	if err := aUUID.UnmarshalBinary(input); err != nil {
		return "", err
	}
	return aUUID.String(), nil
}

// StringUUIDToProtoSnapshotID converts a string UUID to the snapshot ID bytes for protobuf operations.
func StringUUIDToProtoSnapshotID(input string) ([]byte, error) {
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

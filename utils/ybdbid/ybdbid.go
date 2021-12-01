package ybdbid

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

type YBDBID interface {
	Bytes() []byte
	UUID() uuid.UUID
	String() string
}

type defaultYBDBID struct {
	bytes []byte
	uuid  uuid.UUID
	str   string
}

func (id *defaultYBDBID) Bytes() []byte {
	return id.bytes
}

func (id *defaultYBDBID) UUID() uuid.UUID {
	return id.uuid
}

func (id *defaultYBDBID) String() string {
	return id.str
}

// TryParseFromBytes attempts to parse input bytes received from the protobuf
// API as a YugabyteDB ID.
func TryParseFromBytes(input []byte) (YBDBID, error) {
	aUUID := uuid.New()
	if err := aUUID.UnmarshalBinary(input); err != nil {
		return nil, err
	}
	output := &defaultYBDBID{
		bytes: make([]byte, len(input)),
		uuid:  aUUID,
		str:   aUUID.String(),
	}
	copy(output.bytes, input)
	return output, nil
}

// TryParseFromString attempts to parse input string as a YugabyteDB ID.
// Input string can be either a literal UUID or Base64 byte input
// as originally returned by the protobuf API.
func TryParseFromString(input string) (YBDBID, error) {

	// try decoding as base64:
	maybeDecoded, err := base64.StdEncoding.DecodeString(input)
	if err == nil {
		// if succeeded, resulting bytes should be what protobuf would return:
		return TryParseFromBytes(maybeDecoded)
	}

	// if failed, it could be a UUID string, can we parse it as such?
	aUUID, err := uuid.Parse(input)
	if err != nil {
		// no, it's neither base64 encoded, nor looks like UUID:
		return nil, fmt.Errorf("input '%s' is not a valid YugabyteDB ID input", input)
	}

	// it parsed as UUID, we need the bytes too:
	bys, err := aUUID.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("input '%s' is a UUID but could not be marshaled", input)
	}

	output := &defaultYBDBID{
		bytes: make([]byte, len(bys)),
		uuid:  aUUID,
		str:   aUUID.String(),
	}
	copy(output.bytes, bys)
	return output, nil
}

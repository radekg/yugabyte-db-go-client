package ybdbid

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

// This ID type appears in the databasse when a snapshot with incorrect
// configuration is created via RPC. To trigger creating a snapshot
// with such ID, create a YSQL database where:
//  - table namespace is not set
//  - transaction aware is false
var nonUUIDIDTypeR *regexp.Regexp

func init() {
	r, _ := regexp.Compile("^[a-fA-F0-9]{32}$")
	nonUUIDIDTypeR = r
}

// YBDBID represents a parsed YugabyteDB.
type YBDBID interface {
	Bytes() []byte
	String() string
}

type defaultYBDBID struct {
	bytes []byte
	str   string
}

func (id *defaultYBDBID) Bytes() []byte {
	return id.bytes
}

func (id *defaultYBDBID) String() string {
	return id.str
}

// MustParseFromBytes attempts graceful parsing and if there is an error,
// force original input as an ID.
func MustParseFromBytes(input []byte) YBDBID {
	res, err := TryParseFromBytes(input)
	if err != nil {
		if input == nil {
			return &defaultYBDBID{}
		}
		// just return original data:
		output := &defaultYBDBID{
			bytes: make([]byte, len(input)),
			str:   string(input),
		}
		copy(output.bytes, input)
		return output
	}
	return res
}

// TryParseFromBytes attempts to parse input bytes received from the protobuf
// API as a YugabyteDB ID.
func TryParseFromBytes(input []byte) (YBDBID, error) {
	aUUID := uuid.New()
	if err := aUUID.UnmarshalBinary(input); err != nil {

		// if it failed, it could be a third ID type:
		if nonUUIDIDTypeR.Match(input) {
			output := &defaultYBDBID{
				bytes: make([]byte, len(input)),
				str:   string(input),
			}
			copy(output.bytes, input)
			return output, nil
		}

		return nil, err
	}
	output := &defaultYBDBID{
		bytes: make([]byte, len(input)),
		str:   aUUID.String(),
	}
	copy(output.bytes, input)
	return output, nil
}

// MustParseFromString attempts graceful parsing and if there is an error,
// force original input as an ID.
func MustParseFromString(input string) YBDBID {
	res, err := TryParseFromString(input)
	if err != nil {
		// just return original data:
		output := &defaultYBDBID{
			bytes: make([]byte, len(input)),
			str:   string(input),
		}
		copy(output.bytes, []byte(input))
		return output
	}
	return res
}

// TryParseFromString attempts to parse input string as a YugabyteDB ID.
// Input string can be either a literal UUID or Base64 byte input
// as originally returned by the protobuf API.
func TryParseFromString(input string) (YBDBID, error) {

	inputBytes := []byte(input)

	// support third type of ID:
	if nonUUIDIDTypeR.Match(inputBytes) {
		output := &defaultYBDBID{
			bytes: make([]byte, len(inputBytes)),
			str:   input,
		}
		copy(output.bytes, inputBytes)
		return output, nil
	}

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
		str:   aUUID.String(),
	}
	copy(output.bytes, bys)
	return output, nil
}

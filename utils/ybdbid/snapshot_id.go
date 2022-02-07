package ybdbid

/**
	https://yugabyte-db.slack.com/archives/CG0KQF0GG/p1644257955047349?thread_ts=1643868640.696389&cid=CG0KQF0GG
	------------------------------------------------------------------------------------------------------------
	Please interpret the resp.snapshot_id() as an array of 16 bytes. And nothing more. (Not md5/base64/etc.)
	16 bytes. Each byte value is in: 0x00 - 0xFF range.

	In C++ code the decoding only checks the string size and do memcpy:

	Uuid Uuid::TryFullyDecode(const Slice& slice) {
		if (slice.size() != boost::uuids::uuid::static_size()) {
			return Uuid::Nil();
  		}
  		Uuid id;
  		memcpy(id.data(), slice.data(), boost::uuids::uuid::static_size());
  		return id;
	}

	https://yugabyte-db.slack.com/archives/CG0KQF0GG/p1644258454859919?thread_ts=1643868640.696389&cid=CG0KQF0GG
	------------------------------------------------------------------------------------------------------------
	Just because the Snapshot ID for non-transaction-aware-snapshot (old snapshot when
	"transaction_aware=false" - not used now) was passed  through the same PB field. So, in the code
	"32-bytes string" = old non-transactional snapshot UUID as a string. "16 bytes" = new transactional snapshot
	id in binary form. So.. two-in-one.. it's the reason of the complexity. Sorry.
	As the old (non-transactional) snapshots are not used more, you can always expect the 16 bytes.
	The case is only for Snapshot ID.
	Namespace/Table/Tablet ID is a UUID in simple string form. No such complexities.
**/

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// SnapshotID represents a parsed YugabyteDB snapshot ID.
type SnapshotID interface {
	Bytes() []byte
	String() string
	UUID() uuid.UUID
}

type defaultSnapshotID struct {
	bytes []byte
	str   string
	uuuid uuid.UUID
}

func (id *defaultSnapshotID) Bytes() []byte {
	return id.bytes
}

func (id *defaultSnapshotID) String() string {
	return id.str
}

func (id *defaultSnapshotID) UUID() uuid.UUID {
	return id.uuuid
}

// TryParseSnapshotIDFromBytes attempts to parse input bytes received
// from the protobuf API as a YugabyteDB snapshot ID.
func TryParseSnapshotIDFromBytes(input []byte) (SnapshotID, error) {

	if len(input) != 16 {
		return nil, fmt.Errorf("snapshot ID: input must be 16 bytes long")
	}

	aUUID := uuid.New()
	if err := aUUID.UnmarshalBinary(input); err != nil {
		return nil, err
	}
	output := &defaultSnapshotID{
		bytes: make([]byte, len(input)),
		str:   aUUID.String(),
		uuuid: aUUID,
	}
	copy(output.bytes, input)
	return output, nil
}

// TryParseSnapshotIDFromString attempts to parse input string as a YugabyteDB
// snapshot ID. Input string must be a UUIDv4 string.
func TryParseSnapshotIDFromString(input string) (SnapshotID, error) {

	switch len(input) {
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:
	// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9:
		if strings.ToLower(input[:9]) != "urn:uuid:" {
			return nil, fmt.Errorf("snapshot ID: invalid urn prefix: %q", input[:9])
		}
	// {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
	case 36 + 2:
	default:
		return nil, fmt.Errorf("snapshot ID: invalid snapshot ID input")
	}

	// if failed, it could be a UUID string, can we parse it as such?
	aUUID, err := uuid.Parse(input)
	if err != nil {
		// no, it's neither base64 encoded, nor looks like UUID:
		return nil, fmt.Errorf("snapshot ID: input '%s' is not a valid YugabyteDB snapshot ID input", input)
	}

	// it parsed as UUID, we need the bytes too:
	bys, err := aUUID.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("snapshot ID: input '%s' is a UUID but could not be marshaled", input)
	}

	output := &defaultSnapshotID{
		bytes: make([]byte, len(bys)),
		str:   aUUID.String(),
		uuuid: aUUID,
	}
	copy(output.bytes, bys)
	return output, nil
}

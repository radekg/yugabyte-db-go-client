package ybdbid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnapshotIDParsing(t *testing.T) {

	t.Run("it=parses UUID formatted input and back", func(tt *testing.T) {

		validYBDBID := "dfec75ee-290e-4f3b-b965-469a0246c133"
		parsed, err := TryParseSnapshotIDFromString(validYBDBID)
		assert.Nil(tt, err)
		assert.Equal(tt, len(parsed.Bytes()), 16)

		parsedBackViaBytes, err := TryParseSnapshotIDFromBytes(parsed.Bytes())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaBytes.String())

		parsedBackViaString, err := TryParseSnapshotIDFromString(parsed.String())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaString.String())

	})

	t.Run("it=parses UUID formatted wrapped input and back", func(tt *testing.T) {

		validYBDBID := "{dfec75ee-290e-4f3b-b965-469a0246c133}"
		parsed, err := TryParseSnapshotIDFromString(validYBDBID)
		assert.Nil(tt, err)
		assert.Equal(tt, len(parsed.Bytes()), 16)

		parsedBackViaBytes, err := TryParseSnapshotIDFromBytes(parsed.Bytes())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaBytes.String())

		parsedBackViaString, err := TryParseSnapshotIDFromString(parsed.String())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaString.String())

	})

	t.Run("it=parses urn prefixed UUID formatted input and back", func(tt *testing.T) {

		validYBDBID := "urn:uuid:dfec75ee-290e-4f3b-b965-469a0246c133"
		parsed, err := TryParseSnapshotIDFromString(validYBDBID)
		assert.Nil(tt, err)
		assert.Equal(tt, len(parsed.Bytes()), 16)

		parsedBackViaBytes, err := TryParseSnapshotIDFromBytes(parsed.Bytes())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaBytes.String())

		parsedBackViaString, err := TryParseSnapshotIDFromString(parsed.String())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaString.String())

	})

	t.Run("it=handles non-UUID input", func(tt *testing.T) {
		invalidYBDBID := "dfec75ee-290e-4f3b---b965-469a0246c133"
		parsed, err := TryParseSnapshotIDFromString(invalidYBDBID)
		assert.NotNil(tt, err)
		assert.Nil(tt, parsed)
	})

	t.Run("it=handles null byte input", func(tt *testing.T) {
		parsed, err := TryParseFromBytes(nil)
		assert.NotNil(tt, err)
		assert.Nil(tt, parsed)
	})

}

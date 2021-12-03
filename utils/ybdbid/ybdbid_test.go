package ybdbid

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYBDBIDParsing(t *testing.T) {

	t.Run("it=parses base64 encoded input and back", func(tt *testing.T) {

		validBase64YBDBID := "ugsZNctgSFSt0AhKDe7MzA=="
		parsed, err := TryParseFromString(validBase64YBDBID)
		assert.Nil(tt, err)
		assert.Equal(tt, len(parsed.Bytes()), 16)

		parsedBackViaBytes, err := TryParseFromBytes(parsed.Bytes())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaBytes.String())

		parsedBackViaString, err := TryParseFromString(parsed.String())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaString.String())

		assert.Equal(tt, validBase64YBDBID, base64.StdEncoding.EncodeToString(parsedBackViaBytes.Bytes()))
		assert.Equal(tt, validBase64YBDBID, base64.StdEncoding.EncodeToString(parsedBackViaString.Bytes()))

	})

	t.Run("it=parses UUID formatted input and back", func(tt *testing.T) {

		validYBDBID := "dfec75ee-290e-4f3b-b965-469a0246c133"
		parsed, err := TryParseFromString(validYBDBID)
		assert.Nil(tt, err)
		assert.Equal(tt, len(parsed.Bytes()), 16)

		parsedBackViaBytes, err := TryParseFromBytes(parsed.Bytes())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaBytes.String())

		parsedBackViaString, err := TryParseFromString(parsed.String())
		assert.Nil(tt, err)
		assert.Equal(tt, parsed.String(), parsedBackViaString.String())

	})

	t.Run("it=handles invalid base64 encoded input", func(tt *testing.T) {
		invalidBase64YBDBID := base64.StdEncoding.EncodeToString([]byte("this isn't a YBDBID"))
		parsed, err := TryParseFromString(invalidBase64YBDBID)
		assert.NotNil(tt, err)
		assert.Nil(tt, parsed)
	})

	t.Run("it=handles non-UUID input", func(tt *testing.T) {
		invalidYBDBID := "dfec75ee-290e-4f3b---b965-469a0246c133"
		parsed, err := TryParseFromString(invalidYBDBID)
		assert.NotNil(tt, err)
		assert.Nil(tt, parsed)
	})

	t.Run("it=handles null byte input", func(tt *testing.T) {
		parsed, err := TryParseFromBytes(nil)
		assert.NotNil(tt, err)
		assert.Nil(tt, parsed)
	})

	t.Run("it-supports non-UUID ID types", func(tt *testing.T) {
		encodedRepresentation := "ODU4OWYyMzJjYWQxNGJiZWFmNGU4ZTA0NzEwMmYwNmI="
		bys, err := base64.StdEncoding.DecodeString(encodedRepresentation)
		assert.Nil(tt, err)
		parsed, err := TryParseFromBytes(bys)
		assert.Nil(tt, err)
		reparsed, err := TryParseFromString(parsed.String())
		assert.Nil(tt, err)
		reencoded := base64.StdEncoding.EncodeToString([]byte(reparsed.String()))
		assert.Equal(tt, encodedRepresentation, reencoded)
	})

}

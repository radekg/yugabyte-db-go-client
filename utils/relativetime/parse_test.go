package relativetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeOrDuration(t *testing.T) {

	t.Run("it=parses what looks like HT", func(tt *testing.T) {
		f, d, e := ParseTimeOrDuration("1234567890123456")
		assert.Nil(tt, e)
		assert.Equal(tt, time.Duration(0), d)
		assert.Greater(tt, f, uint64(0))
	})

	t.Run("it=parses a valid duration string", func(tt *testing.T) {
		f, d, e := ParseTimeOrDuration("48h15m5s")
		assert.Nil(tt, e)
		assert.Greater(tt, d, time.Duration(0))
		assert.Equal(tt, uint64(0), f)
	})

	t.Run("it=handles invalid strings", func(tt *testing.T) {
		f, d, e := ParseTimeOrDuration("invalid input")
		assert.NotNil(tt, e)
		assert.Equal(tt, time.Duration(0), d)
		assert.Equal(tt, uint64(0), f)
	})

}

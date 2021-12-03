package hybridtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHybridTime(t *testing.T) {

	t.Run("it=handles fixed point in time", func(tt *testing.T) {
		// this is a fixed server clock pointing at exactly '2021-12-03 00:27:34 +0000 UTC':
		fixedPointInTime := uint64(6711260180246810624)
		resultTime := UnixTimeFromHT(fixedPointInTime)
		assert.Equal(tt, "2021-12-03 00:27:34 +0000 UTC", resultTime.UTC().String())
	})

	t.Run("it=converts times correctly", func(tt *testing.T) {
		now := time.Now().UTC()
		resultHT := UnixTimeToHT(now)
		resultTime := UnixTimeFromHT(resultHT)
		assert.Equal(tt, now.Unix(), resultTime.Unix())
	})

	t.Run("it=handles add and substract", func(tt *testing.T) {
		now := time.Now().UTC()
		resultHT := UnixTimeToHT(now)
		afterAddHT := AddDuration(resultHT, time.Duration(time.Hour*48))
		afterSubstractHT := SubstractDuration(afterAddHT, time.Duration(time.Hour*48))
		resultTime := UnixTimeFromHT(afterSubstractHT)
		assert.Equal(tt, now.Unix(), resultTime.Unix())

	})

}

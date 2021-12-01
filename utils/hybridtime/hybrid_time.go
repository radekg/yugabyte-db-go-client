package hybridtime

import "time"

const hybridTimeNumBitsToShift = 12

var hybridTimeLogicalBitsMask = (1 << hybridTimeNumBitsToShift) - 1

// ClockTimestampToHTTimestamp converts the provided timestamp
// to the HybridTime timestamp format. Logical bits are set to 0.
// https://github.com/yugabyte/yugabyte-db/blob/master/java/yb-client/src/main/java/org/yb/util/HybridTimeUtil.java#L55
func ClockTimestampToHTTimestamp(micros uint64) uint64 {
	return micros << hybridTimeNumBitsToShift
}

// HTTimestampToPhysicalAndLogical extracts the physical and logical values
// from an HT timestamp.
// https://github.com/yugabyte/yugabyte-db/blob/master/java/yb-client/src/main/java/org/yb/util/HybridTimeUtil.java#L69
func HTTimestampToPhysicalAndLogical(htTimestamp uint64) []uint64 {
	timestampInMicros := htTimestamp >> hybridTimeNumBitsToShift
	logicalValues := htTimestamp & uint64(hybridTimeLogicalBitsMask)
	return []uint64{timestampInMicros, logicalValues}
}

// PhysicalAndLogicalToHTTimestamp encodes separate physical and logical
// components into a single HT timestamp.
// https://github.com/yugabyte/yugabyte-db/blob/master/java/yb-client/src/main/java/org/yb/util/HybridTimeUtil.java#L82
func PhysicalAndLogicalToHTTimestamp(physical uint64, logical uint64) uint64 {
	return (physical << uint64(hybridTimeNumBitsToShift)) + logical
}

// AddDuration adds a duration to HT.
func AddDuration(ht uint64, d time.Duration) uint64 {
	return ht + ClockTimestampToHTTimestamp(uint64(d.Microseconds()))
}

// SubstractDuration substracts a duration from HT.
func SubstractDuration(ht uint64, d time.Duration) uint64 {
	return ht - ClockTimestampToHTTimestamp(uint64(d.Microseconds()))
}

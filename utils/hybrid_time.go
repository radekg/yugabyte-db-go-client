package utils

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

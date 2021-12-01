package relativetime

import (
	"regexp"
	"strconv"
	"time"
)

// ParseTimeOrDuration parses input string as fixed HT or time.Duration.
func ParseTimeOrDuration(input string) (uint64, time.Duration, error) {

	// The HybridTime is given in microseconds and will contain 16 chars.
	match, _ := regexp.MatchString("^(\\d{16})$", input)
	if match {
		n, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		return n, 0, nil
	}

	d, err := time.ParseDuration(input)
	if err != nil {
		return 0, 0, err
	}

	return 0, d, err

}

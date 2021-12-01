package relativetime

import (
	"time"

	"github.com/radekg/yugabyte-db-go-client/utils/hybridtime"
)

// ServerClockResolver defines a server clock resolver interface.
type ServerClockResolver = func() (uint64, error)

// RelativeOrFixedFuture returns a new restore target in the future if duration is specified,
// or the fixed value.
func RelativeOrFixedFuture(fixed uint64, relative time.Duration, resolver ServerClockResolver) (uint64, error) {
	if relative > 0 {
		serverClock, err := resolver()
		if err != nil {
			return 0, err
		}
		return hybridtime.AddDuration(serverClock, relative), nil
	}
	return fixed, nil
}

// RelativeOrFixedFutureWithFallback returns a new restore target in the future if duration is specified,
// or the fixed value. If both values are 0, resolves new clock value using the resolver.
func RelativeOrFixedFutureWithFallback(fixed uint64, relative time.Duration, resolver ServerClockResolver) (uint64, error) {
	t, err := RelativeOrFixedFuture(fixed, relative, resolver)
	if err != nil {
		return 0, err
	}
	if t > 0 {
		return t, nil
	}
	return resolver()
}

// RelativeOrFixedPast returns a new restore target in the past if duration is specified,
// or the fixed value.
func RelativeOrFixedPast(fixed uint64, relative time.Duration, resolver ServerClockResolver) (uint64, error) {
	if relative > 0 {
		serverClock, err := resolver()
		if err != nil {
			return 0, err
		}
		return hybridtime.SubstractDuration(serverClock, relative), nil
	}
	return fixed, nil
}

// RelativeOrFixedPastWithFallback returns a new restore target in the past if duration is specified,
// or the fixed value. If both values are 0, resolves new clock value using the resolver.
func RelativeOrFixedPastWithFallback(fixed uint64, relative time.Duration, resolver ServerClockResolver) (uint64, error) {
	t, err := RelativeOrFixedPast(fixed, relative, resolver)
	if err != nil {
		return 0, err
	}
	if t > 0 {
		return t, nil
	}
	return resolver()
}

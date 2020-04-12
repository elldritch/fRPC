// Package clock provides an effect abstraction for looking up time.
package clock

import "time"

// A Clock can look up the current time.
type Clock interface {
	Now() time.Time
}

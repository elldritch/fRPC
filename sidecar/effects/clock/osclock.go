package clock

import "time"

var _ Clock = &OSClock{}

// An OSClock provides an implementation of Clock that delegates to the
// underlying syscall.
type OSClock struct{}

// Now delegates to time.Now.
func (*OSClock) Now() time.Time {
	return time.Now()
}

// New constructs an *OSClock.
func New() *OSClock {
	return &OSClock{}
}

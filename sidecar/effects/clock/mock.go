package clock

import "time"

var _ Clock = &Constant{}

// Constant is a mocked clock that returns a specific time.
type Constant struct {
	time time.Time
}

// Now returns the Constant's mocked time.
func (c *Constant) Now() time.Time {
	return c.time
}

// NewConstant returns a new instance of a mock Constant clock.
func NewConstant(t time.Time) *Constant {
	return &Constant{time: t}
}

// Package sensors provides functionality for working with fRPC sensors.
package sensors

// A Tick is a monotically increasing integer, representing the game ticks
// elapsed since the beginning of a specific game.
type Tick uint

// A NetworkID uniquely identifies a circuit network.
type NetworkID int

// A SignalID uniquely identifies a signal by name.
type SignalID string

// A Count specifies the value of a signal.
type Count int

// A Sample contains all circuit network value readings for a specific tick.
type Sample struct {
	Tick     Tick                             `json:"tick"`
	Readings map[NetworkID]map[SignalID]Count `json:"readings"`
}

// A Sensor implements functionality for reading outputs from a running fRPC
// instance.
type Sensor interface {
	// Since returns up to count samples since a specified tick.
	Since(tick Tick, count uint) ([]Sample, error)
	// LatestTick returns the latest tick.
	LatestTick() (Tick, error)
	// PerTick registers a callback which is called on every new reading.
	PerTick(callback func(samples []Sample))
	// DeleteExpired deletes all expired sensor logs.
	DeleteExpired() error
}

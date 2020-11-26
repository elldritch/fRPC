package sensors

import (
	"time"

	"github.com/liftM/fRPC/sidecar/effects/fs"
)

var _ Sensor = &FRPCSensor{}

// A Config contains configuration values for an FRPCSensor.
type Config struct {
	Filesystem fs.Filesystem

	Dir string        // Mod output directory.
	TTL time.Duration // TTL of sensor data.
}

// An FRPCSensor provides a Sensor implementation.
type FRPCSensor struct {
	config Config

	// TODO: sample files are immutable, so we can cache them in-memory to improve
	// lookup speed if needed.
}

// New constructs a new *FRPCSensor.
func New(config Config) *FRPCSensor {
	return &FRPCSensor{config: config}
}

// Since delegates to sensors.Since.
func (s *FRPCSensor) Since(tick Tick, count uint) ([]Sample, error) {
	return Since(s.config.Filesystem, s.config.Dir, tick, count)
}

// LatestTick delegates to sensors.LatestTick.
func (s *FRPCSensor) LatestTick() (Tick, error) {
	return LatestTick(s.config.Filesystem, s.config.Dir)
}

// PerTick registers a callback that runs on every new reading.
func (s *FRPCSensor) PerTick(callback func(samples []Sample)) {
	// TODO: add some way to "unregister" callbacks?
	go func() {
		lastTick, err := s.LatestTick()
		if err != nil {
			panic(err)
		}

		for {
			latestTick, err := s.LatestTick()
			if err != nil {
				panic(err)
			}
			if latestTick != lastTick {
				samples, err := s.Since(lastTick, 1000)
				if err != nil {
					panic(err)
				}
				callback(samples)
				lastTick = latestTick
			}
			time.Sleep(ToDuration(1))
		}
	}()
}

// DeleteExpired calls sensors.Expired to determine the expired samples, and
// then delegates to the underlying filesystem to delete them.
func (s *FRPCSensor) DeleteExpired() error {
	expired, err := Expired(s.config.Filesystem, s.config.Dir, ToTicks(s.config.TTL))
	if err != nil {
		return err
	}
	for _, e := range expired {
		s.config.Filesystem.Delete(e)
	}
	return nil
}

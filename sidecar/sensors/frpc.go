package sensors

import (
	"github.com/liftM/fRPC/sidecar/effects/fs"
)

var _ Sensor = &FRPCSensor{}

// A Config contains configuration values for an FRPCSensor.
type Config struct {
	Filesystem fs.Filesystem

	Dir string // Mod output directory.
	TTL int    // TTL of sensor data in seconds.
}

// An FRPCSensor provides a Sensor implementation.
type FRPCSensor struct {
	config Config
}

// New constructs a new *FRPCSensor.
func New(config Config) *FRPCSensor {
	return &FRPCSensor{config: config}
}

func (s *FRPCSensor) Since(tick Tick, count int) ([]Sample, error) {
	return Since(s.config.Filesystem, s.config.Dir, tick, count)
}

func (s *FRPCSensor) LatestTick() (Tick, error) {
	return LatestTick(s.config.Filesystem, s.config.Dir)
}

func (s *FRPCSensor) DeleteExpired() error {
	expired, err := Expired(s.config.Filesystem, s.config.Dir, s.config.TTL*60)
	if err != nil {
		return err
	}
	for _, e := range expired {
		s.config.Filesystem.Delete(e)
	}
	return nil
}

package sensors

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/liftM/fRPC/sidecar/effects/fs"
)

// TickDuration is the time.Duration of a single Tick.
const TickDuration = time.Second / 60

// Ticks are a measure of a duration of ticks.
type Ticks uint

// ToTicks converts a time.Duration to Ticks.
func ToTicks(d time.Duration) Ticks {
	return Ticks(d / TickDuration)
}

// ToDuration converts a duration of Ticks to a time.Duration.
func ToDuration(t Ticks) time.Duration {
	return time.Duration(time.Duration(t) * TickDuration)
}

// Expired returns a list of paths to files that are more than TTL ticks behind
// the latest tick.
func Expired(fs fs.Filesystem, dir string, ttl Ticks) ([]string, error) {
	files, err := list(fs, dir)
	if err != nil {
		return nil, fmt.Errorf("sensors.Expired: could not read sample directory: %s", err)
	}
	latest := files[len(files)-1]

	var expired []int
	for _, f := range files {
		if latest-f >= int(ttl) {
			expired = append(expired, f)
		}
	}

	var names []string
	for _, e := range expired {
		names = append(names, filepath.Join(dir, strconv.Itoa(e)+".json"))
	}
	return names, nil
}

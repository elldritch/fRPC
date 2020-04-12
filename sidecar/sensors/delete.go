package sensors

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/liftM/fRPC/sidecar/effects/fs"
)

// Expired returns a list of paths to files that are more than TTL ticks behind
// the latest tick.
func Expired(fs fs.Filesystem, dir string, ttl int) ([]string, error) {
	files, err := list(fs, dir)
	if err != nil {
		return nil, fmt.Errorf("sensors.Expired: could not read sample directory: %s", err)
	}
	latest := files[len(files)-1]

	var expired []int
	for _, f := range files {
		if latest-f >= ttl {
			expired = append(expired, f)
		}
	}

	var names []string
	for _, e := range expired {
		names = append(names, filepath.Join(dir, strconv.Itoa(e)+".json"))
	}
	return names, nil
}

package sensors

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/liftM/fRPC/sidecar/effects/fs"
)

// A _Sample is the serialized value logged by fRPC.
type _Sample struct {
	Tick   Tick
	Values []_Value
}

// A _Value is a reading from a sensor for a specific circuit network. Note that
// values may come from the same network multiple times if multiple sensors are
// connected to the same network.
type _Value struct {
	NetworkID NetworkID `json:"network_id"`
	Signals   []_Signal
}

// A _Signal contains both a signal's ID and its count.
type _Signal struct {
	Signal _SignalID
	Count  int
}

// A _SignalID contains both a type and a name.
type _SignalID struct {
	Type string
	Name string
}

// list returns a listing of all ticks with sample files in a directory, sorted
// in ascending order by game tick.
func list(fs fs.Filesystem, dir string) ([]int, error) {
	files, err := fs.List(dir)
	if err != nil {
		return nil, fmt.Errorf("sensors.list: could not read sample directory: %s", err)
	}

	var parsed []int
	for _, f := range files {
		t, err := strconv.Atoi(strings.TrimSuffix(f, ".json"))
		if err != nil {
			return nil, fmt.Errorf("sensors.list: malformed sample file name: %s", err)
		}
		parsed = append(parsed, int(t))
	}

	sort.Ints(parsed)

	return parsed, nil
}

// LatestTick returns the latest logged tick in the directory.
//
// LatestTick should only be called on the mod's output directory, and relies on
// the file naming invariants of that directory.
func LatestTick(fs fs.Filesystem, dir string) (Tick, error) {
	files, err := list(fs, dir)
	if err != nil {
		return 0, fmt.Errorf("sensors.LatestTick: could not read sample directory: %s", err)
	}

	return Tick(files[len(files)-1]), nil
}

// Since reads files in a log directory since a specific game tick, up to a
// maximum count.
//
// Since should only be called on the mod's output directory, and relies on the
// file naming invariants of that directory.
func Since(fs fs.Filesystem, dir string, tick Tick, count uint) ([]Sample, error) {
	files, err := list(fs, dir)
	if err != nil {
		return nil, fmt.Errorf("sensors.Since: could not read sample directory: %s", err)
	}

	// Get the first `count` files since the tick.
	var i uint
	var since []int
	for _, f := range files {
		if f >= int(tick) {
			since = append(since, f)
			i++
			if i >= count {
				break
			}
		}
	}

	// Choose the first `count` number of files.
	var page []string
	for _, t := range since {
		page = append(page, filepath.Join(dir, strconv.Itoa(t)+".json"))
	}

	return ReadFiles(fs, page)
}

// ReadFiles reads a list of files and returns the samples they contain in
// ascending order by tick.
func ReadFiles(fs fs.Filesystem, files []string) ([]Sample, error) {
	var samples []Sample

	for _, f := range files {
		data, err := fs.Read(f)
		if err != nil {
			return nil, fmt.Errorf("sensors.ReadFiles: could not read sample file: %s", err)
		}

		sample, err := Unmarshal(data)
		if err != nil {
			return nil, fmt.Errorf("sensors.ReadFiles: could not unmarshal sample file: %s", err)
		}

		samples = append(samples, sample)
	}

	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Tick < samples[j].Tick
	})

	return samples, nil
}

// Unmarshal unmarshals a single Sample.
func Unmarshal(data []byte) (Sample, error) {
	var sample _Sample
	err := json.Unmarshal(data, &sample)
	if err != nil {
		return Sample{}, fmt.Errorf("sensors.Unmarshal: could not unmarshal sample: %s", err)
	}

	readings := make(map[NetworkID]map[SignalID]Count)
	for _, v := range sample.Values {
		_, ok := readings[v.NetworkID]
		if ok {
			continue
		}
		readings[v.NetworkID] = make(map[SignalID]Count)
		for _, s := range v.Signals {
			readings[v.NetworkID][SignalID(s.Signal.Name)] = Count(s.Count)
		}
	}

	return Sample{
		Tick:     sample.Tick,
		Readings: readings,
	}, nil
}

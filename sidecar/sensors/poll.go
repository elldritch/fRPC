package sensors

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type NetworkID int

type Sample struct {
	Tick   int
	Values []Value
}

type Value struct {
	NetworkID NetworkID `json:"network_id"`
	Signals   []Signal
}

type Signal struct {
	Signal SignalID
	Count  int
}

type SignalID struct {
	Type string
	Name string
}

func (s *Sensors) Poll(dir string) {
	for {
		// List all log files.
		logs, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		// Sort log files in tick order.
		sort.Slice(logs, func(i, j int) bool {
			a, err := strconv.Atoi(strings.TrimSuffix(logs[i].Name(), ".json"))
			if err != nil {
				panic(err)
			}

			b, err := strconv.Atoi(strings.TrimSuffix(logs[j].Name(), ".json"))
			if err != nil {
				panic(err)
			}

			return a < b
		})

		// For each file, record and delete.
		for _, l := range logs {
			bs, err := ioutil.ReadFile(l.Name())
			if err != nil {
				panic(err)
			}

			var sample Sample
			err = json.Unmarshal(bs, &sample)
			if err != nil {
				panic(err)
			}

			s.c <- sample

			err = os.Remove(l.Name())
			if err != nil {
				panic(err)
			}
		}

		// Poll at most once per tick.
		time.Sleep(time.Second / 60)
	}
}

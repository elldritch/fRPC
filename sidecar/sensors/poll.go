package sensors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	lastBucket := 0
	lastPoll := time.Now()

	for {
		logs, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		// Find latest log since last polled tick bucket.
		latestBucket := lastBucket
		for _, info := range logs {
			bucket, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSuffix(info.Name(), ".log"), "frpc_sensors_"))
			if err != nil {
				panic(err)
			}
			if bucket > latestBucket {
				latestBucket = bucket
			}
		}

		// Scan the latest log bucket only if it's changed.
		if latestBucket != lastBucket {
			fmt.Printf("latestBucket: %#v\n", latestBucket)
			f, err := os.Open("frpc_sensors_" + strconv.Itoa(latestBucket) + ".log")
			if err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(f)
			s.mu.Lock()
			for scanner.Scan() {
				var sample Sample
				err := json.Unmarshal(scanner.Bytes(), &sample)
				if err != nil {
					panic(err)
				}

				for _, v := range sample.Values {
					signalMap := make(map[SignalID]int)
					for _, signal := range v.Signals {
						signalMap[signal.Signal] = signal.Count
					}

					s.values[v.NetworkID] = Measurement{
						MeasuredAt: sample.Tick,
						Signals:    signalMap,
					}
				}
			}
			s.mu.Unlock()
		}

		// Poll at most once per second.
		time.Sleep(time.Until(lastPoll.Add(1 * time.Second)))
		lastPoll = time.Now()
		lastBucket = latestBucket
	}
}

package sensors

import (
	"sync"
)

type Measurement struct {
	MeasuredAt int
	Signals    map[SignalID]int
}

type Sensors struct {
	values map[NetworkID]Measurement
	mu     *sync.Mutex
}

func New() *Sensors {
	return &Sensors{
		values: make(map[NetworkID]Measurement),
		mu:     &sync.Mutex{},
	}
}

func (s *Sensors) Read() map[NetworkID]Measurement {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values
}

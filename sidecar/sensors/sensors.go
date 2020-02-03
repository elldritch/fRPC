package sensors

type Sensors struct {
	c chan Sample
}

func New() *Sensors {
	return &Sensors{
		c: make(chan Sample),
	}
}

func (s *Sensors) Subscribe() <-chan Sample {
	return s.c
}

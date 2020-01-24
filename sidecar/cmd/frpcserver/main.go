package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/liftM/fRPC/sidecar/sensors"
)

func main() {
	dir := flag.String("log-dir", "", "directory containing Factorio sensor logs")
	flag.Parse()

	s := sensors.New()
	go s.Poll(*dir)

	for {
		fmt.Printf("%#v\n", s.Read())
		time.Sleep(time.Second)
	}
}

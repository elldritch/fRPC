package main

import (
	"flag"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/liftM/fRPC/sidecar/sensors"
)

func main() {
	pg := flag.String("db", "", "Postgres connection string")
	dir := flag.String("log-dir", "", "directory containing Factorio sensor logs")
	flag.Parse()

	db, err := sqlx.Connect("pgx", *pg)
	if err != nil {
		panic(err)
	}

	s := sensors.New()
	go s.Poll(*dir)

	c := s.Subscribe()
	for {
		sample := <-c
		for _, value := range sample.Values {
			for _, signal := range value.Signals {
				_, err := db.Exec(`
					INSERT INTO circuit_network_signals
						(tick, time, network_id, signal_type, signal_name, count)
					VALUES
						($1, $2, $3, $4, $5, $6)
				`,
					sample.Tick,
					time.Now(),
					value.NetworkID,
					signal.Signal.Type,
					signal.Signal.Name,
					signal.Count,
				)
				if err != nil {
					panic(err)
				}
			}
		}
		fmt.Printf("%#v\n", sample)
	}
}

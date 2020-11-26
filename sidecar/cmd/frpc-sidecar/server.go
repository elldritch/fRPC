package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/liftM/fRPC/sidecar/effects/clock"
	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/liftM/fRPC/sidecar/sensors"
)

type serverConfig struct {
	Clock      clock.Clock
	Filesystem fs.Filesystem

	TTL time.Duration
	Dir string
}

type sensorServer struct {
	sensor sensors.Sensor
	server chi.Router
}

func (s *sensorServer) Start(addr string) error {
	// Set up log rotation.
	go func() {
		for {
			s.sensor.DeleteExpired()
			time.Sleep(sensors.ToDuration(1))
		}
	}()

	// Listen for incoming requests.
	return http.ListenAndServe(addr, s.server)
}

func buildServer(config serverConfig) *sensorServer {
	// Set up services.
	sensor := sensors.New(sensors.Config{
		Filesystem: config.Filesystem,
		TTL:        config.TTL,
		Dir:        config.Dir,
	})

	// Set up router stack.
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/tick", func(w http.ResponseWriter, r *http.Request) {
		now := config.Clock.Now()
		latest, err := sensor.LatestTick()
		if err != nil {
			panic(err)
		}

		res, err := json.Marshal(tickResponse{
			Timestamp: now.Unix(),
			Tick:      latest,
		})
		if err != nil {
			panic(err)
		}

		_, err = w.Write(res)
		if err != nil {
			panic(err)
		}
	})

	r.Get("/samples", func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters and set defaults and restrictions.
		queryparams := r.URL.Query()

		var since sensors.Tick
		sinceQ := queryparams.Get("since")
		if sinceQ != "" {
			// Parse `since`, and constrain it to be non-negative.
			s, err := strconv.Atoi(sinceQ)
			if err != nil {
				data, err := json.Marshal(errorResponse{Error: `"since" must be an integer`})
				if err != nil {
					panic(err)
				}
				w.Write(data)
				return
			}
			if s < 0 {
				data, err := json.Marshal(errorResponse{Error: `"since" must be non-negative`})
				if err != nil {
					panic(err)
				}
				w.Write(data)
				return
			}
			since = sensors.Tick(s)
		} else {
			// Take the last minute of samples.
			latest, err := sensor.LatestTick()
			if err != nil {
				panic(err)
			}
			if latest > 60 {
				since = latest - 60
			} else {
				since = 0
			}
		}

		var count uint
		countQ := queryparams.Get("count")
		if countQ != "" {
			// Parse `count`, and constrain it to be between 1 and 100 (inclusive).
			c, err := strconv.Atoi(countQ)
			if err != nil {
				data, err := json.Marshal(errorResponse{Error: `"count" must be an integer`})
				if err != nil {
					panic(err)
				}
				w.Write(data)
				return
			}
			if c <= 0 {
				data, err := json.Marshal(errorResponse{Error: `"count" must be positive`})
				if err != nil {
					panic(err)
				}
				w.Write(data)
				return
			}
			if c > 100 {
				data, err := json.Marshal(errorResponse{Error: `"count" must be at most 100`})
				if err != nil {
					panic(err)
				}
				w.Write(data)
				return
			}
			count = uint(c)
		} else {
			// By default, return up to 100 samples.
			count = 100
		}

		// TODO: add logger effect - we should log the parsed parameters of every
		// inbound request.

		// Load samples.
		samples, err := sensor.Since(since, count)
		if err != nil {
			panic(err)
		}

		if len(samples) == 0 {
			// All samples are missing.
			data, err := json.Marshal(samplesResponse{
				Missing: &missingInterval{
					Start: since,
				},
				Samples: []sensors.Sample{},
			})
			if err != nil {
				panic(err)
			}
			w.Write(data)
			return
		}

		earliest := samples[0]
		if earliest.Tick-since > 0 {
			// Some samples are missing.
			data, err := json.Marshal(samplesResponse{
				Missing: &missingInterval{
					Start: since,
					End:   earliest.Tick,
				},
				Samples: samples,
			})
			if err != nil {
				panic(err)
			}
			w.Write(data)
			return
		}

		// No samples are missing.
		data, err := json.Marshal(samplesResponse{
			Missing: nil,
			Samples: samples,
		})
		if err != nil {
			panic(err)
		}
		w.Write(data)
	})

	r.Handle("/prometheus", promhttp.Handler())

	return &sensorServer{
		sensor: sensor,
		server: r,
	}
}

type tickResponse struct {
	Timestamp int64        `json:"timestamp"`
	Tick      sensors.Tick `json:"tick"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type samplesResponse struct {
	Missing *missingInterval `json:"missing,omitempty"`
	Samples []sensors.Sample `json:"samples"`
}

type missingInterval struct {
	Start sensors.Tick `json:"start"`
	End   sensors.Tick `json:"end,omitempty"`
}

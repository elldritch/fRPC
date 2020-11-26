package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mitchellh/go-homedir"

	"github.com/liftM/fRPC/sidecar/effects/clock"
	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/liftM/fRPC/sidecar/sensors"
)

func main() {
	// Implement flag parsing.
	defaultDir := ""
	home, err := homedir.Dir()
	if err == nil {
		defaultDir = filepath.Join(home, ".factorio", "script-output")
	}

	addr := flag.String("addr", ":8000", "address for HTTP server to listen on")
	dir := flag.String("dir", defaultDir, "directory containing Factorio sensor logs")
	influxDBURL := flag.String("influx-db-url", "", "URL to InfluxDB instance")
	influxToken := flag.String("influx-db-token", "", "authentication token for InfluxDB")
	influxBucket := flag.String("influx-db-bucket", "", "bucket name for InfluxDB")
	influxOrg := flag.String("influx-db-org", "", "organization name for InfluxDB")
	ttl := flag.Int("ttl", 60, "seconds before deleting sensor data")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, formatParagraphs(helpUsage, formatter{
			Substitute: os.Args[0],
			WrapLength: 80,
		}))
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Implement help command.
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "help":
			if len(os.Args) >= 3 {
				help(os.Args[2])
			} else {
				help("")
			}
		}
	}

	// Implement server.
	server := buildServer(serverConfig{
		Clock:      clock.New(),
		Filesystem: fs.New(),
		TTL:        time.Duration(*ttl) * time.Second,
		Dir:        *dir,
	})

	// Start writing to InfluxDB.
	client := influxdb2.NewClient(*influxDBURL, *influxToken)
	defer client.Close()
	influxWrite := client.WriteAPI(*influxOrg, *influxBucket)
	defer influxWrite.Flush()

	server.sensor.PerTick(func(samples []sensors.Sample) {
		for _, sample := range samples {
			for networkID, signals := range sample.Readings {
				for signalID, value := range signals {
					p := influxdb2.NewPoint("frpc_signal_value", map[string]string{
						"network_id":  strconv.Itoa(int(networkID)),
						"signal_name": string(signalID),
					}, map[string]interface{}{
						"signal_value": int(value),
					}, time.Now())
					influxWrite.WritePoint(p)
				}
			}
		}
	})

	// Start server.
	fmt.Printf("Listening at address %v\n", *addr)
	server.Start(*addr)
}

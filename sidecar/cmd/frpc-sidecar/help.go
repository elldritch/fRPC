package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func help(topic string) {
	switch topic {
	case "api":
		fmt.Fprint(os.Stderr, helpAPI)
		os.Exit(0)
	case "setup":
		fmt.Fprint(os.Stderr, formatParagraphs(helpSetup, formatter{
			Substitute: os.Args[0],
			WrapLength: 80,
		}))
		os.Exit(0)
	case "":
		flag.Usage()
		os.Exit(0)
	default:
		fmt.Fprintln(os.Stderr, `Invalid help topic. Valid help topics: "api".`)
		os.Exit(1)
	}
}

var helpUsage = []string{
	`%s runs an API server that exposes fRPC sensor data.`,
	`It reads sensor logs from a Factorio instance running the fRPC mod, and
	serves those values via an HTTP API. Use "%s help api" to see the API
	documentation. Use "%s help setup" for guidance on connecting this sidecar
	to a running Factorio instance.`,
	`It also cleans up sensor logs when they are older than the configured TTL.
	Sensor logs record all circuit network values on every tick, and get very
	large very quickly. Avoid setting high TTLs unless you have a lot of disk
	space.`,
}

var helpSetup = []string{
	`In order to run %s, install fRPC as a Factorio mod and pass the Factorio
	instance's mod log output directory to %s.`,
	`Most of the time, this directory is located at "~/.factorio/script-output".`,
}

var helpAPI = strings.TrimSpace(strings.ReplaceAll(`
GET /tick

	Returns the latest game tick.

	Example response:

	GET /tick
	{
		// This is the Unix timestamp corresponding to the tick.
		"timestamp": 1586676089,

		// This is the game tick, which measures time from the beginning of the
		// game.
		"tick": 12345,
	}

GET /samples?since=tick&count=integer

	Returns circuit value samples since an integer tick timestamp. Values older
	than the configured TTL are unavailable.

	Example response:

	GET /samples?since=123&count=1
	{
		// The interval of missing values, if any. This will always be a prefix of
		// the requested interval.
		"missing": {
			"start": 123,
			"end": 125
		}

		// The actual samples, up to a maximum of "count".
		"samples": [{
			// The tick of this sample.
			"tick": 126,

			// The values of each circuit network in this sample.
			"values": [{
				"network_id": 6,
				"signals": [
					{
						"signal": {
							"type": "item",
							"name": "copper-ore"
						},
						"count": 4
					}
					// Additional signals on this network...
				]
			}, {
				"network_id": 4,
				"signals": [
					{
						"signal": {
							"type": "item",
							"name": "copper-plate"
						},
						"count": 1
					}
				]
			}
			// Additional circuit networks connected to sensors...
			]
		}
		// Additional samples...
		]
	}
`, "\t", "  ")) + "\n"

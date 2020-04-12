package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
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
		default:
			fmt.Fprintln(os.Stderr, `Invalid subcommand. Valid subcommands: "help".`)
			os.Exit(1)
		}
	}

	// TODO: Implement server.
	_ = dir
	_ = addr
	_ = ttl
}

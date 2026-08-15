// FastNote go_wails3 — Wails v3 application entry point.
//
// Exactly two permitted flags (spec §5.1): --version and --event-file.

package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	eventFile := ""

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--version":
			fmt.Printf("FastNote %s v%s\n", portID, version)
			os.Exit(0)
		case "--event-file":
			if i+1 < len(args) {
				eventFile = args[i+1]
				i++
			}
		default:
			fmt.Fprintf(os.Stderr, "fastnote_wails3: unknown option: %s\n", args[i])
			os.Exit(2)
		}
	}

	runGUI(eventFile)
}

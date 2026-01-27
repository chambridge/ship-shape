// Ship Shape - Test Quality Analysis Tool
// Copyright (c) 2026 Ship Shape Contributors
// Licensed under Apache License 2.0

package main

import (
	"fmt"
	"os"
)

// Version information (set via ldflags during build)
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

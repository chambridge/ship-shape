// Ship Shape - Version Command
// Copyright (c) 2026 Ship Shape Contributors
// Licensed under Apache License 2.0

package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	shortVersion bool
	jsonVersion  bool
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Display version information for Ship Shape.

Examples:
  # Show full version info
  shipshape version

  # Show version number only
  shipshape version --short

  # Output as JSON
  shipshape version --json`,
	RunE: runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "show version number only")
	versionCmd.Flags().BoolVar(&jsonVersion, "json", false, "output as JSON")
}

type versionInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

func runVersion(_ *cobra.Command, _ []string) error {
	info := versionInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}

	if shortVersion {
		fmt.Println(Version)
		return nil
	}

	if jsonVersion {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}

		fmt.Println(string(data))

		return nil
	}

	// Standard output
	fmt.Printf("Ship Shape v%s\n", info.Version)
	fmt.Printf("Build Date: %s\n", info.BuildTime)
	fmt.Printf("Git Commit: %s\n", info.GitCommit)
	fmt.Printf("Go Version: %s\n", info.GoVersion)
	fmt.Printf("Platform: %s\n", info.Platform)

	return nil
}

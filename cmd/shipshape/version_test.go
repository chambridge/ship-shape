package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/chambridge/ship-shape/internal/testutil"
	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Save original version info
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit

	// Set test values
	Version = "1.0.0"
	BuildTime = "2026-01-01T00:00:00Z"
	GitCommit = "abc123"

	// Restore original values
	t.Cleanup(func() {
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	})

	tests := []struct {
		name           string
		args           []string
		wantContains   []string
		wantNotContain []string
		wantErr        bool
	}{
		{
			name: "default output",
			args: []string{"version"},
			wantContains: []string{
				"Ship Shape v1.0.0",
				"Build Date: 2026-01-01T00:00:00Z",
				"Git Commit: abc123",
				"Go Version:",
				"Platform:",
			},
		},
		{
			name: "short output",
			args: []string{"version", "--short"},
			wantContains: []string{
				"1.0.0",
			},
			wantNotContain: []string{
				"Ship Shape",
				"Build Date",
				"Git Commit",
			},
		},
		{
			name: "json output",
			args: []string{"version", "--json"},
			wantContains: []string{
				`"version": "1.0.0"`,
				`"build_time": "2026-01-01T00:00:00Z"`,
				`"git_commit": "abc123"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for each test
			shortVersion = false
			jsonVersion = false

			// Capture output
			stdout, _ := testutil.CaptureOutput(t, func() {
				// Create a new root command for isolation
				cmd := &cobra.Command{Use: "shipshape"}
				cmd.AddCommand(versionCmd)

				// Set args
				cmd.SetArgs(tt.args)

				// Execute
				err := cmd.Execute()

				// Check error
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			// Check contains
			for _, want := range tt.wantContains {
				if !strings.Contains(stdout, want) {
					t.Errorf("Output does not contain %q:\n%s", want, stdout)
				}
			}

			// Check not contains
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(stdout, notWant) {
					t.Errorf("Output should not contain %q:\n%s", notWant, stdout)
				}
			}
		})
	}
}

func TestVersionCommandJSON(t *testing.T) {
	// Save original version info
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit

	// Set test values
	Version = "1.0.0"
	BuildTime = "2026-01-01T00:00:00Z"
	GitCommit = "abc123"

	// Restore original values
	t.Cleanup(func() {
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	})

	// Reset flags
	shortVersion = false
	jsonVersion = false

	// Create a new root command for isolation
	cmd := &cobra.Command{Use: "shipshape"}
	cmd.AddCommand(versionCmd)

	// Capture output
	stdout, _ := testutil.CaptureOutput(t, func() {
		cmd.SetArgs([]string{"version", "--json"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	// Parse JSON
	var info versionInfo
	if err := json.Unmarshal([]byte(stdout), &info); err != nil {
		t.Fatalf("Failed to parse JSON: %v\nOutput: %s", err, stdout)
	}

	// Verify fields
	if info.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", info.Version, "1.0.0")
	}

	if info.BuildTime != "2026-01-01T00:00:00Z" {
		t.Errorf("BuildTime = %q, want %q", info.BuildTime, "2026-01-01T00:00:00Z")
	}

	if info.GitCommit != "abc123" {
		t.Errorf("GitCommit = %q, want %q", info.GitCommit, "abc123")
	}

	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}

	if info.Platform == "" {
		t.Error("Platform should not be empty")
	}
}

func TestRunVersion(t *testing.T) {
	// Save original version info
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit

	// Set test values
	Version = "test-version"
	BuildTime = "test-time"
	GitCommit = "test-commit"

	// Restore original values
	t.Cleanup(func() {
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	})

	tests := []struct {
		name         string
		wantContains string
		shortFlag    bool
		jsonFlag     bool
	}{
		{
			name:         "default format",
			shortFlag:    false,
			jsonFlag:     false,
			wantContains: "Ship Shape v",
		},
		{
			name:         "short format",
			shortFlag:    true,
			jsonFlag:     false,
			wantContains: "test-version",
		},
		{
			name:         "json format",
			shortFlag:    false,
			jsonFlag:     true,
			wantContains: `"version"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set flags
			shortVersion = tt.shortFlag
			jsonVersion = tt.jsonFlag

			// Capture output
			stdout, _ := testutil.CaptureOutput(t, func() {
				if err := runVersion(nil, nil); err != nil {
					t.Errorf("runVersion() error = %v", err)
				}
			})

			if !strings.Contains(stdout, tt.wantContains) {
				t.Errorf("Output does not contain %q:\n%s", tt.wantContains, stdout)
			}
		})
	}
}

func TestVersionInfo(t *testing.T) {
	// Test versionInfo struct marshaling
	info := versionInfo{
		Version:   "1.0.0",
		BuildTime: "2026-01-01",
		GitCommit: "abc123",
		GoVersion: "go1.21.0",
		Platform:  "darwin/arm64",
	}

	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal versionInfo: %v", err)
	}

	var decoded versionInfo
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal versionInfo: %v", err)
	}

	if decoded != info {
		t.Errorf("Decoded versionInfo does not match original:\ngot  %+v\nwant %+v", decoded, info)
	}
}

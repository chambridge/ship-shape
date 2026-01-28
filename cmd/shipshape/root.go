// Ship Shape - Root Command
// Copyright (c) 2026 Ship Shape Contributors
// Licensed under Apache License 2.0

package main

import (
	"fmt"
	"os"

	"github.com/chambridge/ship-shape/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	quiet   bool
	noColor bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shipshape",
	Short: "Ship Shape - Test Quality Analysis Tool",
	Long: `Ship Shape analyzes your codebase to assess test quality and coverage.

It provides comprehensive insights into:
  • Test coverage metrics (line, branch, function)
  • Test quality and smells detection
  • Performance and execution speed
  • Testing tools and best practices adoption
  • Code maintainability and organization

Ship Shape supports Go, Python, JavaScript/TypeScript, Java, and Rust.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: .shipshape.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (errors only)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")

	// Bind flags to viper
	// nolint:errcheck // BindPFlag always succeeds for valid flags
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	// nolint:errcheck
	_ = viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	// nolint:errcheck
	_ = viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory, repository root, and home directory
		viper.AddConfigPath(".")
		viper.AddConfigPath(findRepositoryRoot())
		viper.AddConfigPath("$HOME/.config/shipshape")
		viper.AddConfigPath("/etc/shipshape")

		// Config file name (without extension)
		viper.SetConfigName(".shipshape")
		viper.SetConfigType("yaml")
	}

	// Environment variable prefix
	viper.SetEnvPrefix("SHIPSHAPE")
	viper.AutomaticEnv()

	// Read config file
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; using defaults
			logger.Debug("No config file found, using defaults")

			return
		}

		// Config file found but another error occurred
		logger.Error("Error reading config file", "error", err)
		os.Exit(1)
	}

	logger.Debug("Using config file", "path", viper.ConfigFileUsed())
}

// findRepositoryRoot searches for the repository root by looking for .git directory
func findRepositoryRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(fmt.Sprintf("%s/.git", dir)); err == nil {
			return dir
		}

		parent := fmt.Sprintf("%s/..", dir)
		if parent == dir {
			// Reached filesystem root
			break
		}

		dir = parent
	}

	return "."
}

// initLogger initializes the logger based on CLI flags.
func initLogger() {
	// Determine log level based on flags
	var level logger.Level

	switch {
	case quiet:
		level = logger.LevelError
	case verbose:
		level = logger.LevelDebug
	default:
		level = logger.LevelInfo
	}

	// Create logger configuration
	cfg := logger.Config{
		Level:   level,
		Format:  "text",
		Output:  os.Stderr,
		NoColor: noColor,
	}

	// Initialize and set as default
	l := logger.New(cfg)
	logger.SetDefault(l)

	// Log initialization message at debug level
	logger.Debug("Logger initialized",
		"level", level,
		"format", cfg.Format,
		"no_color", noColor,
	)
}

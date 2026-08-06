// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package engine

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tforceaio/llm-db/config"
	"github.com/tforceaio/llm-db/schema/dotenv"
)

var version = "0.1.0"

// Log error message
func logProgramError(logger zerolog.Logger, err error) {
	if err != nil {
		logger.Err(err).Msg("Unexpected error has occurred. Program will exit.")
	}
}

// InitApp initializes the logger.
func InitApp() zerolog.Logger {
	logger := config.InitZerolog()

	fmt.Printf("LLM-DB Utiltity v%s.\nCopyright (C) %d T-Force I/O.\nLicensed under MIT license.\n\n", version, 2025)

	return logger
}

// Execute runs the CLI application.
func Execute() {
	rootCmd := &cobra.Command{
		Use:     "llm-db",
		Short:   fmt.Sprintf("LLM DB Utiltity v%s\n", version),
		Long:    fmt.Sprintf("LLM-DB Utiltity v%s.\nCopyright (C) %d T-Force I/O.\nLicensed under MIT license.", version, 2025),
		Version: version,
	}

	rootCmd.AddCommand(ExportCmd())
	rootCmd.AddCommand(FormatCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Load environement variables from envFile, or from .env file in current directory.
// Return a boolean indicating whether the environment variables were successfully loaded.
func loadEnvFile(envFile string, override bool) bool {
	if envFile == "" {
		if _, err := os.Stat(".env"); err != nil {
			return false
		}
		if err := dotenv.LoadDotEnv(".env", override); err != nil {
			return false
		}
	} else {
		if err := dotenv.LoadDotEnv(envFile, override); err != nil {
			return false
		}
	}
	return true
}

// Update variable with the value of the environment variable if it exists.
func resolveEnvVar(envKey string, variable *string) {
	envVar := os.Getenv(envKey)
	if envVar != "" {
		*variable = envVar
	}
}

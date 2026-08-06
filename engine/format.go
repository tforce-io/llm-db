// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tforceaio/llm-db/schema/opencode_cloud"
)

const (
	opencodeModelsPath = "json/refs/opencode/models.json"
)

// FormatModule handles reformatting of refs JSON files.
type FormatModule struct {
	logger zerolog.Logger
}

// Return a new FormatModule.
func NewFormatModule(logger zerolog.Logger, cmdName string) *FormatModule {
	return &FormatModule{
		logger: logger.With().Str("module", "format").Str("cmd", cmdName).Logger(),
	}
}

// Reformat opencode models JSON file (sort keys alphabetically, pretty-print with 2-space indent).
func (m *FormatModule) OpenCodeModels(inputPath, outputPath string) error {
	m.logger.Info().Str("path", inputPath).Msg("Loading models...")
	modelsData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", inputPath, err)
	}

	models, err := opencode_cloud.LoadModels(modelsData)
	if err != nil {
		return fmt.Errorf("failed to parse models: %w", err)
	}

	m.logger.Info().Int("count", len(models)).Msg("Loaded models.")

	m.logger.Info().Str("path", outputPath).Msg("Writing formatted models...")
	outputData, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal models: %w", err)
	}

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	absPath, _ := filepath.Abs(outputPath)
	m.logger.Info().Str("path", absPath).Msg("Format complete.")
	return nil
}

// Decorator to log error occurred when calling handlers.
func (m *FormatModule) logError(err error) {
	logProgramError(m.logger, err)
}

// Define Cobra Command for Format module.
func FormatCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "format",
		Short: "Format refs JSON files.",
	}

	opencodeModelsCmd := &cobra.Command{
		Use:   "opencode-models",
		Short: "Reformat opencode models.json (sort keys, pretty-print).",
		Run: func(cmd *cobra.Command, args []string) {
			logger := InitApp()
			flags := ParseFormatFlags(cmd, args)
			m := NewFormatModule(logger, "opencode-models")
			m.logError(m.OpenCodeModels(flags.Input, flags.Output))
		},
	}

	opencodeModelsCmd.Flags().String("input", opencodeModelsPath, "Input JSON file path")
	opencodeModelsCmd.Flags().String("output", opencodeModelsPath, "Output JSON file path")

	rootCmd.AddCommand(opencodeModelsCmd)
	return rootCmd
}

// FormatFlags contains all flags used by FormatModule.
type FormatFlags struct {
	Input  string
	Output string
}

// Extract all flags from a Cobra Command.
func ParseFormatFlags(cmd *cobra.Command, args []string) *FormatFlags {
	input, _ := cmd.Flags().GetString("input")
	output, _ := cmd.Flags().GetString("output")

	return &FormatFlags{
		Input:  input,
		Output: output,
	}
}

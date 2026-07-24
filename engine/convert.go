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
	"github.com/tforceaio/llm-db/schema/llmdb"
	"github.com/tforceaio/llm-db/schema/opencode"
)

const (
	modelsPath    = "json/models.json"
	providersPath = "json/providers.json"
	outputDir     = "exported"
	outputPath    = "exported/opencode.json"
)

// ConvertModule handles conversion of llmdb models to different formats.
type ConvertModule struct {
	logger zerolog.Logger
}

// Return a new ConvertModule.
func NewConvertModule(logger zerolog.Logger, cmdName string) *ConvertModule {
	return &ConvertModule{
		logger: logger.With().Str("module", "convert").Str("cmd", cmdName).Logger(),
	}
}

// Export LLMDB models into OpenCode config.
func (m *ConvertModule) OpenCode(apiKey string) error {
	m.logger.Info().Msg("Loading models...")
	modelsData, err := os.ReadFile(modelsPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", modelsPath, err)
	}

	models, err := llmdb.LoadModels(modelsData)
	if err != nil {
		return fmt.Errorf("failed to parse models: %w", err)
	}

	m.logger.Info().Int("count", len(models.Models)).Msg("Loaded models.")

	m.logger.Info().Msg("Loading providers...")
	providersData, err := os.ReadFile(providersPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", providersPath, err)
	}

	providers, err := llmdb.LoadProviders(providersData)
	if err != nil {
		return fmt.Errorf("failed to parse providers: %w", err)
	}

	m.logger.Info().Int("count", len(providers.Providers)).Msg("Loaded providers.")

	cfg := m.buildOpenCodeConfig(models, providers, apiKey)

	m.logger.Info().Str("dir", outputDir).Msg("Creating output directory...")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", outputDir, err)
	}

	m.logger.Info().Str("path", outputPath).Msg("Writing opencode config...")
	outputData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}

	absPath, _ := filepath.Abs(outputPath)
	m.logger.Info().Str("path", absPath).Msg("Conversion complete.")
	return nil
}

func (m *ConvertModule) buildOpenCodeConfig(models *llmdb.Models, providers *llmdb.Providers, apiKey string) *opencode.RootConfig {
	cfg := &opencode.RootConfig{
		Schema:       "https://opencode.ai/config.json",
		SmallModel:   "opencode/deepseek-v4-flash-free",
		DefaultAgent: "plan",
		Agents: map[string]opencode.AgentConfigs{
			"build": {Color: "primary"},
			"plan":  {Color: "secondary"},
		},
		Providers: make(map[string]opencode.ProviderConfigs),
	}

	ollamaProvider := opencode.ProviderConfigs{Models: make(map[string]opencode.ModelConfig)}
	openCodeProvider := opencode.ProviderConfigs{Models: make(map[string]opencode.ModelConfig)}
	bifrostProvider := opencode.ProviderConfigs{Models: make(map[string]opencode.ModelConfig)}

	for _, model := range models.Models {
		for deployKey, deployment := range model.Deployments {
			mc := m.buildOpenCodeModel(model, deployment)

			if deployKey == "opencode-zen" {
				openCodeProvider.Models[deployment.ID] = opencode.ModelConfig{
					Name: mc.Name,
					Cost: mc.Cost,
				}
			}
			if deployKey == "ollama" {
				ollamaProvider.Models[deployment.ID] = mc
			} else {
				modelID := fmt.Sprintf("%s/%s", deployKey, deployment.ID)

				providerName := ""
				if p, ok := providers.Providers[deployKey]; ok {
					providerName = p.Name
				}
				mc.ID = modelID
				mc.Name = fmt.Sprintf("%s (%s)", model.Name, providerName)

				bifrostProvider.Models[modelID] = mc
			}
		}
	}

	if ollamaP, ok := providers.Providers["ollama"]; ok {
		ollamaProvider.Name = ollamaP.Name
		ollamaProvider.NPM = ollamaP.NPM
		opts := &opencode.ProviderOptions{}
		if ollamaP.URI != "" {
			opts.BaseUrl = ollamaP.URI
		}
		if apiKey != "" {
			opts.ApiKey = apiKey
		}
		if opts.BaseUrl != "" || opts.ApiKey != "" {
			ollamaProvider.Options = opts
		}
	}

	if bifrostP, ok := providers.Providers["bifrost"]; ok {
		bifrostProvider.Name = bifrostP.Name
		bifrostProvider.NPM = bifrostP.NPM
		opts := &opencode.ProviderOptions{}
		if bifrostP.URI != "" {
			opts.BaseUrl = bifrostP.URI
		}
		if apiKey != "" {
			opts.ApiKey = apiKey
		}
		if opts.BaseUrl != "" || opts.ApiKey != "" {
			bifrostProvider.Options = opts
		}
	}

	cfg.Providers["ollama"] = ollamaProvider
	cfg.Providers["opencode"] = openCodeProvider
	cfg.Providers["bifrost"] = bifrostProvider

	return cfg
}

func (m *ConvertModule) buildOpenCodeModel(model llmdb.Model, deployment llmdb.Deployment) opencode.ModelConfig {
	reasoning := containsString(model.Capabilities, llmdb.CapabilityReasoning)
	temperature := containsString(model.Capabilities, llmdb.CapabilityTemperature)
	toolCall := containsString(model.Capabilities, llmdb.CapabilityTools) || containsString(model.Capabilities, llmdb.CapabilityFunctionCall)

	mc := opencode.ModelConfig{
		ID:                   deployment.ID,
		Name:                 model.Name,
		Cost:                 &opencode.ModelCost{Input: model.Cost.Input, Output: model.Cost.Output},
		Limit:                m.resolveLimit(model.Limit, deployment.Limit),
		Modalities:           &opencode.ModelModalities{Input: model.Modalities.Input, Output: model.Modalities.Output},
		SupportReasoning:     &reasoning,
		SupportTemperature:   &temperature,
		SupportsFunctionCall: &toolCall,
	}

	return mc
}

func (m *ConvertModule) resolveLimit(base llmdb.ModelLimit, override *llmdb.ModelLimit) *opencode.ModelLimit {
	if override != nil {
		return &opencode.ModelLimit{Context: override.Context, Output: override.Output}
	}
	return &opencode.ModelLimit{Context: base.Context, Output: base.Output}
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Decorator to log error occurred when calling handlers.
func (m *ConvertModule) logError(err error) {
	logProgramError(m.logger, err)
}

// Define Cobra Command for Convert opencode module.
func ConvertCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert models to different formats.",
	}

	var apiKey string

	opencodeCmd := &cobra.Command{
		Use:   "opencode",
		Short: "Export LLMDB models into OpenCode config.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := InitApp()
			m := NewConvertModule(logger, "opencode")
			m.logError(m.OpenCode(apiKey))
		},
	}

	opencodeCmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for ollama and bifrost providers")

	rootCmd.AddCommand(opencodeCmd)
	return rootCmd
}

// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tforceaio/llm-db/schema/bifrost"
	"github.com/tforceaio/llm-db/schema/llmdb"
	"github.com/tforceaio/llm-db/schema/opencode"
)

const (
	modelsPath         = "json/models.json"
	providersPath      = "json/providers.json"
	outputDir          = "exported"
	bifrostOutputPath  = "exported/bifrost-models.json"
	openCodeOutputPath = "exported/opencode.json"
)

// ExportModule handles conversion of llmdb models to different formats.
type ExportModule struct {
	logger zerolog.Logger
}

// Return a new ExportModule.
func NewExportModule(logger zerolog.Logger, cmdName string) *ExportModule {
	return &ExportModule{
		logger: logger.With().Str("module", "export").Str("cmd", cmdName).Logger(),
	}
}

// Export LLMDB models into Bifrost models JSON.
func (m *ExportModule) Bifrost() error {
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

	bModels := m.buildBifrostModels(models)

	m.logger.Info().Str("dir", outputDir).Msg("Creating output directory...")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", outputDir, err)
	}

	m.logger.Info().Str("path", bifrostOutputPath).Msg("Writing bifrost models...")
	outputData, err := json.MarshalIndent(bModels, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal models: %w", err)
	}

	if err := os.WriteFile(bifrostOutputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", bifrostOutputPath, err)
	}

	absPath, _ := filepath.Abs(bifrostOutputPath)
	m.logger.Info().Str("path", absPath).Msg("Export complete.")
	return nil
}

func (m *ExportModule) buildBifrostModels(models *llmdb.Models) bifrost.Models {
	result := make(bifrost.Models)

	for _, model := range models.Models {
		for deployKey, deployment := range model.Deployments {
			key := deployment.ID
			if deployKey != "ollama" {
				key = deployKey + "/" + deployment.ID
			}

			cost := m.resolveCost(model.Cost, deployment.Cost)
			limit := m.resolveLimit(model.Limit, deployment.Limit)

			mode := "chat"
			if len(model.Capabilities) == 0 {
				mode = "embedding"
			}

			bModel := bifrost.Model{
				Provider:                    deployKey,
				BaseModel:                   deployment.ID,
				Mode:                        mode,
				InputCost:                   bifrost.Float64(divMillionth(cost.Input)),
				OutputCost:                  bifrost.Float64(divMillionth(cost.Output)),
				CacheReadInputTokenCost:     bifrost.Float64(divMillionth(cost.Cache)),
				CacheCreationInputTokenCost: bifrost.Float64(divMillionth(cost.CacheWrite)),
				MaxInputTokens:              limit.Context,
				MaxOutputTokens:             limit.Output,
				MaxTokens:                   limit.Context,
			}

			if containsString(model.Capabilities, llmdb.CapabilityFunctionCall) ||
				containsString(model.Capabilities, llmdb.CapabilityTools) {
				bModel.SupportsFunctionCall = boolPtr(true)
			}
			if containsString(model.Capabilities, llmdb.CapabilityReasoning) {
				bModel.SupportsReasoning = boolPtr(true)
			}
			if containsString(model.Capabilities, llmdb.CapabilityStructured) {
				bModel.SupportsStructured = boolPtr(true)
			}
			if containsString(model.Capabilities, llmdb.CapabilityToolChoice) {
				bModel.SupportsToolChoice = boolPtr(true)
			}
			if containsString(model.Capabilities, llmdb.CapabilityVision) {
				bModel.SupportsVision = boolPtr(true)
			}

			result[key] = bModel
		}
	}
	return result
}

// Export LLMDB models into OpenCode config.
func (m *ExportModule) OpenCode(ollamaURL, ollamaApiKey, bifrostURL, bifrostApiKey string, updateGlobalConfig bool) error {
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

	cfg := m.buildOpenCodeConfig(models, providers, ollamaURL, ollamaApiKey, bifrostURL, bifrostApiKey)

	m.logger.Info().Str("dir", outputDir).Msg("Creating output directory...")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", outputDir, err)
	}

	m.logger.Info().Str("path", openCodeOutputPath).Msg("Writing opencode config...")
	outputData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(openCodeOutputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", openCodeOutputPath, err)
	}

	absPath, _ := filepath.Abs(openCodeOutputPath)
	m.logger.Info().Str("path", absPath).Msg("Export complete.")

	if updateGlobalConfig {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home dir: %w", err)
		}

		configDir := filepath.Join(home, ".config", "opencode")
		globalConfigPath := filepath.Join(configDir, "opencode.json")

		m.logger.Info().Str("dir", configDir).Msg("Creating global config directory...")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", configDir, err)
		}

		m.logger.Info().Str("path", globalConfigPath).Msg("Writing global opencode config...")
		if err := os.WriteFile(globalConfigPath, outputData, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", globalConfigPath, err)
		}

		m.logger.Info().Str("path", globalConfigPath).Msg("Global config written.")
	}

	return nil
}

func (m *ExportModule) buildOpenCodeConfig(models *llmdb.Models, providers *llmdb.Providers, ollamaURL, ollamaApiKey, bifrostURL, bifrostApiKey string) *opencode.RootConfig {
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
	bifrostProvider := opencode.ProviderConfigs{Models: make(map[string]opencode.ModelConfig)}

	for _, model := range models.Models {
		for deployKey, deployment := range model.Deployments {
			mc := m.buildOpenCodeModel(model, deployment)
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
		if ollamaURL != "" {
			opts.BaseUrl = ollamaURL
		}
		if ollamaApiKey != "" {
			opts.ApiKey = ollamaApiKey
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
		if bifrostURL != "" {
			opts.BaseUrl = bifrostURL
		}
		if bifrostApiKey != "" {
			opts.ApiKey = bifrostApiKey
		}
		if opts.BaseUrl != "" || opts.ApiKey != "" {
			bifrostProvider.Options = opts
		}
	}

	cfg.Providers["ollama"] = ollamaProvider
	cfg.Providers["bifrost"] = bifrostProvider

	return cfg
}

func (m *ExportModule) buildOpenCodeModel(model llmdb.Model, deployment llmdb.Deployment) opencode.ModelConfig {
	reasoning := containsString(model.Capabilities, llmdb.CapabilityReasoning)
	temperature := containsString(model.Capabilities, llmdb.CapabilityTemperature)
	toolCall := containsString(model.Capabilities, llmdb.CapabilityTools) || containsString(model.Capabilities, llmdb.CapabilityFunctionCall)

	cost := m.resolveCost(model.Cost, deployment.Cost)
	limit := m.resolveLimit(model.Limit, deployment.Limit)

	mc := opencode.ModelConfig{
		ID:                   deployment.ID,
		Name:                 model.Name,
		Cost:                 &opencode.ModelCost{Input: cost.Input, Output: cost.Output, CacheRead: cost.Cache, CacheWrite: cost.CacheWrite},
		Limit:                &opencode.ModelLimit{Context: limit.Context, Output: limit.Output},
		Modalities:           &opencode.ModelModalities{Input: model.Modalities.Input, Output: model.Modalities.Output},
		SupportReasoning:     &reasoning,
		SupportTemperature:   &temperature,
		SupportsFunctionCall: &toolCall,
	}

	return mc
}

func (m *ExportModule) resolveCost(base llmdb.ModelCost, override *llmdb.ModelCost) llmdb.ModelCost {
	if override != nil {
		return *override
	}
	return base
}

func (m *ExportModule) resolveLimit(base llmdb.ModelLimit, override *llmdb.ModelLimit) *llmdb.ModelLimit {
	if override != nil {
		return override
	}
	return &base
}

// Decorator to log error occurred when calling handlers.
func (m *ExportModule) logError(err error) {
	logProgramError(m.logger, err)
}

// Define Cobra Command for Export opencode module.
func ExportCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "export",
		Short: "Export models to different formats.",
	}

	bifrostCmd := &cobra.Command{
		Use:   "bifrost",
		Short: "Export LLMDB models into Bifrost models JSON.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := InitApp()
			m := NewExportModule(logger, "bifrost")
			m.logError(m.Bifrost())
		},
	}
	rootCmd.AddCommand(bifrostCmd)

	var bifrostApiKey, bifrostURL, ollamaApiKey, ollamaURL string
	var updateConfig bool

	opencodeCmd := &cobra.Command{
		Use:   "opencode",
		Short: "Export LLMDB models into OpenCode config.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := InitApp()
			m := NewExportModule(logger, "opencode")
			m.logError(m.OpenCode(ollamaURL, ollamaApiKey, bifrostURL, bifrostApiKey, updateConfig))
		},
	}

	opencodeCmd.Flags().StringVar(&bifrostApiKey, "bifrost-api-key", "", "API key for Bifrost provider")
	opencodeCmd.Flags().StringVar(&bifrostURL, "bifrost-url", "", "Override Bifrost provider base URL")
	opencodeCmd.Flags().StringVar(&ollamaApiKey, "ollama-api-key", "", "API key for Ollama provider")
	opencodeCmd.Flags().StringVar(&ollamaURL, "ollama-url", "", "Override Ollama provider base URL")
	opencodeCmd.Flags().BoolVarP(&updateConfig, "deploy", "u", false, "Update global OpenCode config file")

	rootCmd.AddCommand(opencodeCmd)
	return rootCmd
}

func boolPtr(b bool) *bool {
	return &b
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Divides f by 1,000,000 and strips floating-point artifacts by
// parsing the result through a 10-significant-figure string representation.
func divMillionth(f float64) float64 {
	s := strconv.FormatFloat(f/1_000_000, 'g', 12, 64)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

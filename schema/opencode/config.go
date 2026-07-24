// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode

import "encoding/json"

type RootConfig struct {
	Schema       string                     `json:"$schema,omitempty"`
	SmallModel   string                     `json:"small_model,omitempty"`
	DefaultAgent string                     `json:"default_agent,omitempty"`
	Agents       map[string]AgentConfigs    `json:"agent,omitempty"`
	MCPs         map[string]MCPConfigs      `json:"mcp,omitempty"`
	Providers    map[string]ProviderConfigs `json:"provider,omitempty"`
}

type AgentConfigs struct {
	Color string `json:"color,omitempty"`
}

type MCPConfigs struct {
	Type    string            `json:"type,omitempty"`
	URL     string            `json:"url,omitempty"`
	Enabled bool              `json:"enabled,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

type ProviderConfigs struct {
	Name    string                 `json:"name,omitempty"`
	NPM     string                 `json:"npm,omitempty"`
	Options *ProviderOptions       `json:"options,omitempty"`
	Models  map[string]ModelConfig `json:"models,omitempty"`
}

type ProviderOptions struct {
	BaseUrl string `json:"baseURL,omitempty"`
	ApiKey  string `json:"apiKey,omitempty"`
}

type ModelConfig struct {
	ID                   string           `json:"id,omitempty"`
	Name                 string           `json:"name,omitempty"`
	Description          string           `json:"description,omitempty"`
	Cost                 *ModelCost       `json:"cost,omitempty"`
	Limit                *ModelLimit      `json:"limit,omitempty"`
	Modalities           *ModelModalities `json:"modalities,omitempty"`
	SupportReasoning     *bool            `json:"reasoning,omitempty"`
	SupportTemperature   *bool            `json:"temperature,omitempty"`
	SupportsFunctionCall *bool            `json:"tool_call,omitempty"`
}

type ModelCost struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type ModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type ModelModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

func LoadConfig(data []byte) (*RootConfig, error) {
	var cfg RootConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

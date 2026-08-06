// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode_cloud

import "encoding/json"

type API map[string]Provider

type Provider struct {
	ID     string              `json:"id"`
	Env    []string            `json:"env"`
	NPM    string              `json:"npm"`
	API    string              `json:"api"`
	Name   string              `json:"name"`
	Doc    string              `json:"doc"`
	Models map[string]ApiModel `json:"models"`
}

type ApiModel struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Family           string              `json:"family"`
	Attachment       bool                `json:"attachment"`
	Reasoning        bool                `json:"reasoning"`
	ReasoningOptions []ReasoningOption   `json:"reasoning_options"`
	ToolCall         bool                `json:"tool_call"`
	StructuredOutput bool                `json:"structured_output"`
	Temperature      bool                `json:"temperature"`
	Knowledge        string              `json:"knowledge"`
	ReleaseDate      string              `json:"release_date"`
	LastUpdated      string              `json:"last_updated"`
	OpenWeights      bool                `json:"open_weights"`
	Status           string              `json:"status"`
	Modalities       *ProviderModalities `json:"modalities"`
	Limit            *ProviderLimit      `json:"limit"`
	Cost             *Cost               `json:"cost"`
	Interleaved      *Interleaved        `json:"interleaved,omitempty"`
}

type ReasoningOption struct {
	Type   string   `json:"type"`
	Values []string `json:"values,omitempty"`
	Min    *int     `json:"min,omitempty"`
	Max    *int     `json:"max,omitempty"`
}

type ProviderModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type ProviderLimit struct {
	Context int `json:"context"`
	Input   int `json:"input,omitempty"`
	Output  int `json:"output"`
}

type Cost struct {
	Input           float64    `json:"input"`
	Output          float64    `json:"output"`
	CacheRead       float64    `json:"cache_read,omitempty"`
	CacheWrite      float64    `json:"cache_write,omitempty"`
	Tiers           []CostTier `json:"tiers,omitempty"`
	ContextOver200k *Cost      `json:"context_over_200k,omitempty"`
}

type CostTier struct {
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	CacheRead  float64 `json:"cache_read,omitempty"`
	CacheWrite float64 `json:"cache_write,omitempty"`
	Tier       Tier    `json:"tier"`
}

type Tier struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Interleaved struct {
	Field  string `json:"field"`
	IsBool bool   `json:"-"`
	Value  bool   `json:"-"`
}

func (i *Interleaved) UnmarshalJSON(data []byte) error {
	if string(data) == "true" {
		i.IsBool = true
		i.Value = true
		return nil
	}
	if string(data) == "false" {
		i.IsBool = true
		i.Value = false
		return nil
	}
	type Alias Interleaved
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(i),
	}
	return json.Unmarshal(data, aux)
}

func LoadAPI(data []byte) (API, error) {
	var api API
	if err := json.Unmarshal(data, &api); err != nil {
		return nil, err
	}
	return api, nil
}

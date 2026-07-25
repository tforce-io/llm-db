// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package bifrost

import (
	"encoding/json"
	"strconv"
)

// Float64 is a float64 that marshals to decimal notation (e.g. 0.0000001) instead of scientific notation (e.g. 1e-7).
type Float64 float64

func (f Float64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	*f = Float64(v)
	return nil
}

type Models map[string]Model

type Model struct {
	Provider                    string  `json:"provider"`
	BaseModel                   string  `json:"base_model"`
	Mode                        string  `json:"mode"`
	InputCost                   Float64 `json:"input_cost_per_token"`
	OutputCost                  Float64 `json:"output_cost_per_token"`
	CacheReadInputTokenCost     Float64 `json:"cache_read_input_token_cost,omitempty"`
	CacheCreationInputTokenCost Float64 `json:"cache_creation_input_token_cost,omitempty"`
	MaxInputTokens              int     `json:"max_input_tokens"`
	MaxOutputTokens             int     `json:"max_output_tokens"`
	MaxTokens                   int     `json:"max_tokens"`
	SupportsFunctionCall        *bool   `json:"supports_function_calling,omitempty"`
	SupportsReasoning           *bool   `json:"supports_reasoning,omitempty"`
	SupportsStructured          *bool   `json:"supports_response_schema,omitempty"`
	SupportsToolChoice          *bool   `json:"supports_tool_choice,omitempty"`
	SupportsVision              *bool   `json:"supports_vision,omitempty"`
}

func LoadModels(data []byte) (Models, error) {
	var models Models
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return models, nil
}

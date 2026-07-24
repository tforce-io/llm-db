// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package bifrost

import "encoding/json"

type Models map[string]Model

type Model struct {
	Provider             string  `json:"provider"`
	BaseModel            string  `json:"base_model"`
	Mode                 string  `json:"mode"`
	InputCost            float64 `json:"input_cost_per_token"`
	OutputCost           float64 `json:"output_cost_per_token"`
	MaxInputTokens       int     `json:"max_input_tokens"`
	MaxOutputTokens      int     `json:"max_output_tokens"`
	MaxTokens            int     `json:"max_tokens"`
	SupportsFunctionCall *bool   `json:"supports_function_calling,omitempty"`
	SupportsReasoning    *bool   `json:"supports_reasoning,omitempty"`
	SupportsStructured   *bool   `json:"supports_response_schema,omitempty"`
	SupportsToolChoice   *bool   `json:"supports_tool_choice,omitempty"`
	SupportsVision       *bool   `json:"supports_vision,omitempty"`
}

func LoadModels(data []byte) (Models, error) {
	var models Models
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return models, nil
}

// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package bifrost

import (
	"encoding/json"

	"github.com/tforceaio/llm-db/common"
)

type Models map[string]Model

type Model struct {
	Provider  string `json:"provider,omitempty"`
	BaseModel string `json:"base_model,omitempty"`
	Mode      string `json:"mode,omitempty"`

	TextInputCost    *common.Float64 `json:"input_cost_per_token,omitempty"`
	TextOutputCost   *common.Float64 `json:"output_cost_per_token,omitempty"`
	CacheReadCost    *common.Float64 `json:"cache_read_input_token_cost,omitempty"`
	CacheWriteCost   *common.Float64 `json:"cache_creation_input_token_cost,omitempty"`
	VisionInputCost  *common.Float64 `json:"input_cost_per_image_token,omitempty"`
	VisionOutputCost *common.Float64 `json:"output_cost_per_image_token,omitempty"`
	ImageInputCost   *common.Float64 `json:"input_cost_per_image,omitempty"`
	ImageOutputCost  *common.Float64 `json:"output_cost_per_image,omitempty"`
	PixelInputCost   *common.Float64 `json:"input_cost_per_pixel,omitempty"`
	PixelOutputCost  *common.Float64 `json:"output_cost_per_pixel,omitempty"`

	MaxInputTokens  *int `json:"max_input_tokens,omitempty"`
	MaxOutputTokens *int `json:"max_output_tokens,omitempty"`
	MaxTokens       *int `json:"max_tokens,omitempty"`

	SupportsFunctionCall *bool `json:"supports_function_calling,omitempty"`
	SupportsReasoning    *bool `json:"supports_reasoning,omitempty"`
	SupportsStructured   *bool `json:"supports_response_schema,omitempty"`
	SupportsToolChoice   *bool `json:"supports_tool_choice,omitempty"`
	SupportsVision       *bool `json:"supports_vision,omitempty"`
}

func LoadModels(data []byte) (Models, error) {
	var models Models
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return models, nil
}

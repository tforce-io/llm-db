// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package llmdb

import (
	"encoding/json"

	"github.com/tforceaio/llm-db/common"
)

const (
	CapabilityEmbedding       = "embedding"
	CapabilityFunctionCall    = "function_call"
	CapabilityImageEdit       = "image_edit"
	CapabilityImageGeneration = "image_generation"
	CapabilityReasoning       = "reasoning"
	CapabilityResponseFormat  = "response_format"
	CapabilityStructured      = "structured"
	CapabilityTemperature     = "temperature"
	CapabilityToolChoice      = "tool_choice"
	CapabilityTools           = "tools"
	CapabilityVision          = "vision"
)

var ValidCapabilities = map[string]bool{
	CapabilityEmbedding:       true,
	CapabilityFunctionCall:    true,
	CapabilityImageEdit:       true,
	CapabilityImageGeneration: true,
	CapabilityReasoning:       true,
	CapabilityResponseFormat:  true,
	CapabilityStructured:      true,
	CapabilityTemperature:     true,
	CapabilityToolChoice:      true,
	CapabilityTools:           true,
	CapabilityVision:          true,
}

const (
	ModalityText  = "text"
	ModalityPdf   = "pdf"
	ModalityImage = "image"
	ModalityVideo = "video"
	ModalityAudio = "audio"
)

var ValidModalities = map[string]bool{
	ModalityText:  true,
	ModalityPdf:   true,
	ModalityImage: true,
	ModalityVideo: true,
	ModalityAudio: true,
}

type Models struct {
	Schema string `json:"$schema,omitempty"`
	Models map[string]Model
}

func (m *Models) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if schemaRaw, ok := raw["$schema"]; ok {
		if err := json.Unmarshal(schemaRaw, &m.Schema); err != nil {
			return err
		}
		delete(raw, "$schema")
	}

	m.Models = make(map[string]Model)
	for key, value := range raw {
		var model Model
		if err := json.Unmarshal(value, &model); err != nil {
			return err
		}
		m.Models[key] = model
	}

	return nil
}

func (m Models) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	if m.Schema != "" {
		result["$schema"] = m.Schema
	}
	for key, model := range m.Models {
		result[key] = model
	}
	return json.Marshal(result)
}

type Model struct {
	Name         string                `json:"name"`
	Base         string                `json:"base,omitempty"`
	Home         string                `json:"home"`
	OSS          string                `json:"oss,omitempty"`
	Specs        string                `json:"specs,omitempty"`
	Dev          string                `json:"dev,omitempty"`
	Capabilities []string              `json:"capabilities"`
	Cost         ModelCost             `json:"cost"`
	Limit        ModelLimit            `json:"limit"`
	Modalities   ModelModalities       `json:"modalities"`
	Deployments  map[string]Deployment `json:"deployments"`
}

type ModelCost struct {
	TextInput    *common.Float64 `json:"input,omitempty"`
	TextOutput   *common.Float64 `json:"output,omitempty"`
	CacheRead    *common.Float64 `json:"cache_read,omitempty"`
	CacheWrite   *common.Float64 `json:"cache_write,omitempty"`
	VisionInput  *common.Float64 `json:"vision_input,omitempty"`
	VisionOutput *common.Float64 `json:"vision_output,omitempty"`
	ImageInput   *common.Float64 `json:"image_input,omitempty"`
	ImageOutput  *common.Float64 `json:"image_output,omitempty"`
	PixelInput   *common.Float64 `json:"pixel_input,omitempty"`
	PixelOutput  *common.Float64 `json:"pixel_output,omitempty"`
}

type ModelLimit struct {
	Context *int `json:"context,omitempty"`
	Output  *int `json:"output,omitempty"`
}

type ModelModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type Deployment struct {
	ID           string      `json:"id"`
	Limit        *ModelLimit `json:"limit,omitempty"`
	Cost         *ModelCost  `json:"cost,omitempty"`
	Capabilities []string    `json:"capabilities,omitempty"`
}

func LoadModels(data []byte) (*Models, error) {
	var models Models
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return &models, nil
}

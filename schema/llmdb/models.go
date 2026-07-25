// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package llmdb

import "encoding/json"

const (
	CapabilityFunctionCall   = "function_call"
	CapabilityReasoning      = "reasoning"
	CapabilityResponseFormat = "response_format"
	CapabilityStructured     = "structured"
	CapabilityTemperature    = "temperature"
	CapabilityToolChoice     = "tool_choice"
	CapabilityTools          = "tools"
	CapabilityVision         = "vision"
)

var ValidCapabilities = map[string]bool{
	CapabilityFunctionCall:   true,
	CapabilityReasoning:      true,
	CapabilityResponseFormat: true,
	CapabilityStructured:     true,
	CapabilityTemperature:    true,
	CapabilityToolChoice:     true,
	CapabilityTools:          true,
	CapabilityVision:         true,
}

const (
	ModalityText  = "text"
	ModalityImage = "image"
	ModalityVideo = "video"
	ModalityAudio = "audio"
)

var ValidModalities = map[string]bool{
	ModalityText:  true,
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
	Input      float64 `json:"input"`
	Output     float64 `json:"output"`
	Cache      float64 `json:"cache,omitempty"`
	CacheWrite float64 `json:"cache_write,omitempty"`
}

type ModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
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

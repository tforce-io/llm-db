// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode_cloud

import "encoding/json"

type Models map[string]Model

type Model struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	Family           string      `json:"family"`
	Attachment       bool        `json:"attachment"`
	Reasoning        bool        `json:"reasoning"`
	ToolCall         bool        `json:"tool_call"`
	StructuredOutput bool        `json:"structured_output"`
	Temperature      bool        `json:"temperature"`
	Knowledge        string      `json:"knowledge"`
	ReleaseDate      string      `json:"release_date"`
	LastUpdated      string      `json:"last_updated"`
	OpenWeights      bool        `json:"open_weights"`
	License          string      `json:"license"`
	Modalities       *Modalities `json:"modalities"`
	Limit            *Limit      `json:"limit"`
	Weights          []Weight    `json:"weights"`
	Benchmarks       []Benchmark `json:"benchmarks"`
	Links            []Link      `json:"links"`
}

type Modalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type Limit struct {
	Context int `json:"context"`
	Input   int `json:"input,omitempty"`
	Output  int `json:"output"`
}

type Weight struct {
	Label        string `json:"label"`
	URL          string `json:"url"`
	Quantization string `json:"quantization,omitempty"`
}

type Benchmark struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Metric  string  `json:"metric"`
	Source  string  `json:"source"`
	Harness string  `json:"harness,omitempty"`
	Dataset string  `json:"dataset,omitempty"`
	Date    string  `json:"date,omitempty"`
	Variant string  `json:"variant,omitempty"`
	Version string  `json:"version,omitempty"`
}

type Link struct {
	Label string `json:"label"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

func LoadModels(data []byte) (Models, error) {
	var models Models
	if err := json.Unmarshal(data, &models); err != nil {
		return nil, err
	}
	return models, nil
}

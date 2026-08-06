// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode_cloud

import (
	"os"
	"testing"
)

func TestLoadModels(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/models.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	models, err := LoadModels(data)
	if err != nil {
		t.Fatalf("failed to load models: %v", err)
	}

	if len(models) == 0 {
		t.Error("expected non-empty models")
	}

	if _, ok := models["zhipuai/glm-5"]; !ok {
		t.Error("expected zhipuai/glm-5 model")
	}

	if _, ok := models["zhipuai/glm-4.6v"]; !ok {
		t.Error("expected zhipuai/glm-4.6v model")
	}
}

func TestModelEntryFields(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/models.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	models, err := LoadModels(data)
	if err != nil {
		t.Fatalf("failed to load models: %v", err)
	}

	model := models["zhipuai/glm-4.6v"]

	if model.ID != "zhipuai/glm-4.6v" {
		t.Errorf("expected ID 'zhipuai/glm-4.6v', got '%s'", model.ID)
	}

	if model.Name != "GLM-4.6V" {
		t.Errorf("expected name 'GLM-4.6V', got '%s'", model.Name)
	}

	if model.Family != "glm" {
		t.Errorf("expected family 'glm', got '%s'", model.Family)
	}

	if !model.Reasoning {
		t.Error("expected reasoning to be true")
	}

	if !model.ToolCall {
		t.Error("expected tool_call to be true")
	}

	if !model.OpenWeights {
		t.Error("expected open_weights to be true")
	}

	if model.Modalities == nil {
		t.Fatal("expected modalities to be set")
	}

	if len(model.Modalities.Input) == 0 {
		t.Error("expected input modalities")
	}

	if model.Limit == nil {
		t.Fatal("expected limit to be set")
	}

	if model.Limit.Context != 128000 {
		t.Errorf("expected context limit 128000, got %d", model.Limit.Context)
	}
}

func TestWeights(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/models.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	models, err := LoadModels(data)
	if err != nil {
		t.Fatalf("failed to load models: %v", err)
	}

	model := models["zhipuai/glm-4.6v"]

	if len(model.Weights) == 0 {
		t.Fatal("expected weights to be set")
	}

	weight := model.Weights[0]
	if weight.Label != "Hugging Face" {
		t.Errorf("expected label 'Hugging Face', got '%s'", weight.Label)
	}

	if weight.URL == "" {
		t.Error("expected URL to be set")
	}
}

func TestBenchmarks(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/models.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	models, err := LoadModels(data)
	if err != nil {
		t.Fatalf("failed to load models: %v", err)
	}

	model := models["zhipuai/glm-5"]

	if len(model.Benchmarks) == 0 {
		t.Fatal("expected benchmarks to be set")
	}

	benchmark := model.Benchmarks[0]
	if benchmark.Name != "SWE-Bench Verified" {
		t.Errorf("expected name 'SWE-Bench Verified', got '%s'", benchmark.Name)
	}

	if benchmark.Score == 0 {
		t.Error("expected score to be set")
	}

	if benchmark.Metric == "" {
		t.Error("expected metric to be set")
	}

	if benchmark.Source == "" {
		t.Error("expected source to be set")
	}
}

func TestStructuredOutput(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/models.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	models, err := LoadModels(data)
	if err != nil {
		t.Fatalf("failed to load models: %v", err)
	}

	model := models["zhipuai/glm-5.1"]

	if !model.StructuredOutput {
		t.Error("expected structured_output to be true")
	}
}

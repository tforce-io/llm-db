// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode_cloud

import (
	"os"
	"testing"
)

func TestLoadAPI(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	if len(api) == 0 {
		t.Error("expected non-empty API")
	}

	if _, ok := api["zhipuai"]; !ok {
		t.Error("expected zhipuai provider")
	}

	if _, ok := api["lucidquery"]; !ok {
		t.Error("expected lucidquery provider")
	}

	if _, ok := api["anyapi"]; !ok {
		t.Error("expected anyapi provider")
	}
}

func TestProviderModels(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	provider, ok := api["zhipuai"]
	if !ok {
		t.Fatal("zhipuai provider not found")
	}

	if provider.ID != "zhipuai" {
		t.Errorf("expected provider ID 'zhipuai', got '%s'", provider.ID)
	}

	if provider.Name != "Zhipu AI" {
		t.Errorf("expected provider name 'Zhipu AI', got '%s'", provider.Name)
	}

	if len(provider.Models) == 0 {
		t.Error("expected zhipuai to have models")
	}

	if _, ok := provider.Models["glm-5"]; !ok {
		t.Error("expected glm-5 model in zhipuai")
	}
}

func TestModelFields(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	provider := api["zhipuai"]
	model := provider.Models["glm-5"]

	if model.ID != "glm-5" {
		t.Errorf("expected model ID 'glm-5', got '%s'", model.ID)
	}

	if model.Name != "GLM-5" {
		t.Errorf("expected model name 'GLM-5', got '%s'", model.Name)
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

	if model.Modalities == nil {
		t.Fatal("expected modalities to be set")
	}

	if len(model.Modalities.Input) == 0 {
		t.Error("expected input modalities")
	}

	if model.Limit == nil {
		t.Fatal("expected limit to be set")
	}

	if model.Limit.Context != 204800 {
		t.Errorf("expected context limit 204800, got %d", model.Limit.Context)
	}

	if model.Cost == nil {
		t.Fatal("expected cost to be set")
	}

	if model.Cost.Input != 1 {
		t.Errorf("expected input cost 1, got %f", model.Cost.Input)
	}
}

func TestReasoningOptions(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	provider := api["zhipuai"]
	model := provider.Models["glm-5.2"]

	if len(model.ReasoningOptions) == 0 {
		t.Fatal("expected reasoning options")
	}

	opt := model.ReasoningOptions[0]
	if opt.Type != "effort" {
		t.Errorf("expected type 'effort', got '%s'", opt.Type)
	}

	if len(opt.Values) == 0 {
		t.Error("expected values for effort type")
	}
}

func TestCostTiers(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	provider := api["crossmodel"]
	model := provider.Models["gemini/gemini-3.1-pro-preview"]

	if model.Cost == nil {
		t.Fatal("expected cost to be set")
	}

	if len(model.Cost.Tiers) == 0 {
		t.Error("expected cost tiers")
	}

	if model.Cost.ContextOver200k == nil {
		t.Error("expected context_over_200k to be set")
	}
}

func TestInterleaved(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/opencode/api.json")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	api, err := LoadAPI(data)
	if err != nil {
		t.Fatalf("failed to load API: %v", err)
	}

	provider := api["zhipuai"]
	model := provider.Models["glm-5"]

	if model.Interleaved == nil {
		t.Fatal("expected interleaved to be set")
	}

	if model.Interleaved.Field != "reasoning_content" {
		t.Errorf("expected field 'reasoning_content', got '%s'", model.Interleaved.Field)
	}
}

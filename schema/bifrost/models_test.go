// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package bifrost_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tforceaio/llm-db/schema/bifrost"
)

func TestLoadModels(t *testing.T) {
	data, err := os.ReadFile("../../json/samples/bifrost.json")
	if err != nil {
		t.Fatalf("failed to read sample file: %v", err)
	}

	models, err := bifrost.LoadModels(data)
	if err != nil {
		t.Fatalf("LoadModels failed: %v", err)
	}

	if len(models) != 3 {
		t.Fatalf("expected 3 models, got %d", len(models))
	}

	glm := models["glm-4.7-flash:30b-a3b"]
	if glm.Provider != "ollama" {
		t.Errorf("expected provider 'ollama', got '%s'", glm.Provider)
	}
	if glm.BaseModel != "glm-4.7-flash:30b-a3b" {
		t.Errorf("expected base_model 'glm-4.7-flash:30b-a3b', got '%s'", glm.BaseModel)
	}
	if glm.Mode != "chat" {
		t.Errorf("expected mode 'chat', got '%s'", glm.Mode)
	}
	if glm.TextInputCost == nil || *glm.TextInputCost != 0.000000060 {
		t.Errorf("expected input_cost_per_token 0.000000060, got %v", glm.TextInputCost)
	}
	if glm.TextOutputCost == nil || *glm.TextOutputCost != 0.000000400 {
		t.Errorf("expected output_cost_per_token 0.000000400, got %v", glm.TextOutputCost)
	}
	if glm.MaxInputTokens == nil || *glm.MaxInputTokens != 202752 {
		t.Errorf("expected max_input_tokens 202752, got %v", glm.MaxInputTokens)
	}
	if glm.MaxOutputTokens == nil || *glm.MaxOutputTokens != 131072 {
		t.Errorf("expected max_output_tokens 131072, got %v", glm.MaxOutputTokens)
	}
	if glm.SupportsFunctionCall == nil || !*glm.SupportsFunctionCall {
		t.Error("expected supports_function_calling to be true")
	}

	qwen := models["qwen3.6:27b"]
	if qwen.SupportsVision == nil || !*qwen.SupportsVision {
		t.Error("expected qwen3.6:27b supports_vision to be true")
	}
	if qwen.CacheReadCost == nil || *qwen.CacheReadCost != 0.000000250 {
		t.Errorf("expected cache_read_input_token_cost 0.000000250, got %v", qwen.CacheReadCost)
	}
	if qwen.CacheWriteCost == nil || *qwen.CacheWriteCost != 0.000004000 {
		t.Errorf("expected cache_creation_input_token_cost 0.000004000, got %v", qwen.CacheWriteCost)
	}

	pickle := models["opencode-zen/big-pickle"]
	if pickle.Provider != "opencode-zen" {
		t.Errorf("expected provider 'opencode-zen', got '%s'", pickle.Provider)
	}
	if pickle.SupportsVision == nil || *pickle.SupportsVision {
		t.Error("expected opencode-zen/big-pickle supports_vision to be false")
	}
}

func TestLoadImageModels(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/bifrost/datasheet.json")
	if err != nil {
		t.Fatalf("failed to read datasheet file: %v", err)
	}

	models, err := bifrost.LoadModels(data)
	if err != nil {
		t.Fatalf("LoadModels failed: %v", err)
	}

	flux := models["aiml/flux/schnell"]
	if flux.Mode != "image_generation" {
		t.Errorf("expected mode 'image_generation', got '%s'", flux.Mode)
	}
	if flux.ImageOutputCost == nil || *flux.ImageOutputCost != 0.004 {
		t.Errorf("expected output_cost_per_image 0.004, got %v", flux.ImageOutputCost)
	}
	if flux.TextInputCost != nil {
		t.Error("expected input_cost_per_token to be absent")
	}
	if flux.MaxInputTokens != nil {
		t.Error("expected max_input_tokens to be absent")
	}

	dalle := models["256-x-256/dall-e-2"]
	if dalle.PixelInputCost == nil || *dalle.PixelInputCost != 2.4414e-7 {
		t.Errorf("expected input_cost_per_pixel 2.4414e-7, got %v", dalle.PixelInputCost)
	}
	if dalle.PixelOutputCost == nil || *dalle.PixelOutputCost != 0 {
		t.Errorf("expected output_cost_per_pixel 0, got %v", dalle.PixelOutputCost)
	}
	if dalle.ImageOutputCost != nil {
		t.Error("expected output_cost_per_image to be absent")
	}

	canvas := models["1024-x-1024/50-steps/bedrock/amazon.nova-canvas-v1:0"]
	if canvas.Provider != "bedrock" {
		t.Errorf("expected provider 'bedrock', got '%s'", canvas.Provider)
	}
	if canvas.ImageOutputCost == nil || *canvas.ImageOutputCost != 0.06 {
		t.Errorf("expected output_cost_per_image 0.06, got %v", canvas.ImageOutputCost)
	}
	if canvas.MaxInputTokens == nil || *canvas.MaxInputTokens != 2600 {
		t.Errorf("expected max_input_tokens 2600, got %v", canvas.MaxInputTokens)
	}
}

func TestRoundTrip(t *testing.T) {
	data, err := os.ReadFile("../../json/samples/bifrost.json")
	if err != nil {
		t.Fatalf("failed to read sample file: %v", err)
	}

	models, err := bifrost.LoadModels(data)
	if err != nil {
		t.Fatalf("LoadModels failed: %v", err)
	}

	marshaled, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var original map[string]interface{}
	if err := json.Unmarshal(data, &original); err != nil {
		t.Fatalf("failed to unmarshal original: %v", err)
	}

	var roundtripped map[string]interface{}
	if err := json.Unmarshal(marshaled, &roundtripped); err != nil {
		t.Fatalf("failed to unmarshal roundtripped: %v", err)
	}

	assertJSONEqual(t, original, roundtripped, "original", "roundtripped")
}

func assertJSONEqual(t *testing.T, expected, actual interface{}, expName, actName string) {
	t.Helper()

	if fmt.Sprintf("%T", expected) != fmt.Sprintf("%T", actual) &&
		!(isNil(expected) && isNil(actual)) {
		t.Errorf("type mismatch at %s vs %s: expected %T, got %T", expName, actName, expected, actual)
		return
	}

	if isNil(expected) && isNil(actual) {
		return
	}

	switch e := expected.(type) {
	case map[string]interface{}:
		a := actual.(map[string]interface{})
		for k, v := range e {
			av, exists := a[k]
			if !exists {
				t.Errorf("missing key '%s' in %s", k, actName)
				continue
			}
			assertJSONEqual(t, v, av, fmt.Sprintf("%s.%s", expName, k), fmt.Sprintf("%s.%s", actName, k))
		}
		for k := range a {
			if _, exists := e[k]; !exists {
				t.Errorf("unexpected key '%s' in %s", k, actName)
			}
		}
	case []interface{}:
		a := actual.([]interface{})
		if len(e) != len(a) {
			t.Errorf("%s has %d elements, %s has %d", expName, len(e), actName, len(a))
			return
		}
		for i := range e {
			assertJSONEqual(t, e[i], a[i], fmt.Sprintf("%s[%d]", expName, i), fmt.Sprintf("%s[%d]", actName, i))
		}
	case float64:
		if expected != actual {
			t.Errorf("mismatch at %s vs %s: expected %v, got %v", expName, actName, expected, actual)
		}
	case string:
		if expected != actual {
			t.Errorf("mismatch at %s vs %s: expected %q, got %q", expName, actName, expected, actual)
		}
	case bool:
		if expected != actual {
			t.Errorf("mismatch at %s vs %s: expected %v, got %v", expName, actName, expected, actual)
		}
	}
}

func TestImageModelsRoundTrip(t *testing.T) {
	data, err := os.ReadFile("../../json/refs/bifrost/datasheet.json")
	if err != nil {
		t.Fatalf("failed to read datasheet file: %v", err)
	}

	models, err := bifrost.LoadModels(data)
	if err != nil {
		t.Fatalf("LoadModels failed: %v", err)
	}

	modeledKeys := map[string]bool{
		"provider":   true,
		"base_model": true,
		"mode":       true,

		"input_cost_per_token":            true,
		"output_cost_per_token":           true,
		"cache_read_input_token_cost":     true,
		"cache_creation_input_token_cost": true,
		"input_cost_per_image_token":      true,
		"output_cost_per_image_token":     true,
		"input_cost_per_image":            true,
		"output_cost_per_image":           true,
		"input_cost_per_pixel":            true,
		"output_cost_per_pixel":           true,

		"max_input_tokens":  true,
		"max_output_tokens": true,
		"max_tokens":        true,

		"supports_function_calling": true,
		"supports_reasoning":        true,
		"supports_response_schema":  true,
		"supports_tool_choice":      true,
		"supports_vision":           true,
	}

	var original map[string]map[string]interface{}
	if err := json.Unmarshal(data, &original); err != nil {
		t.Fatalf("failed to unmarshal original: %v", err)
	}

	for key, model := range models {
		marshaled, err := json.Marshal(model)
		if err != nil {
			t.Fatalf("failed to marshal %s: %v", key, err)
		}

		var roundtripped map[string]interface{}
		if err := json.Unmarshal(marshaled, &roundtripped); err != nil {
			t.Fatalf("failed to unmarshal roundtripped %s: %v", key, err)
		}

		for field := range roundtripped {
			if !modeledKeys[field] {
				t.Errorf("unexpected key '%s' in roundtripped %s", field, key)
			}
			if _, exists := original[key][field]; !exists {
				t.Errorf("key '%s' present in roundtripped %s but absent in original", field, key)
			}
		}
		for field, value := range original[key] {
			if !modeledKeys[field] {
				continue
			}
			// Explicit empty strings normalize to absent on marshal (omitempty).
			if s, ok := value.(string); ok && s == "" {
				continue
			}
			if _, exists := roundtripped[field]; !exists {
				t.Errorf("modeled key '%s' present in original %s but missing in roundtripped", field, key)
				continue
			}
			assertJSONEqual(t, map[string]interface{}{field: value}, map[string]interface{}{field: roundtripped[field]}, key, key)
		}
	}
}

func isNil(v interface{}) bool {
	return v == nil
}

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
	if glm.InputCost != 0.000000060 {
		t.Errorf("expected input_cost_per_token 0.000000060, got %f", glm.InputCost)
	}
	if glm.OutputCost != 0.000000400 {
		t.Errorf("expected output_cost_per_token 0.000000400, got %f", glm.OutputCost)
	}
	if glm.MaxInputTokens != 202752 {
		t.Errorf("expected max_input_tokens 202752, got %d", glm.MaxInputTokens)
	}
	if glm.MaxOutputTokens != 131072 {
		t.Errorf("expected max_output_tokens 131072, got %d", glm.MaxOutputTokens)
	}
	if glm.SupportsFunctionCall == nil || !*glm.SupportsFunctionCall {
		t.Error("expected supports_function_calling to be true")
	}

	qwen := models["qwen3.6:27b"]
	if qwen.SupportsVision == nil || !*qwen.SupportsVision {
		t.Error("expected qwen3.6:27b supports_vision to be true")
	}

	pickle := models["opencode-zen/big-pickle"]
	if pickle.Provider != "opencode-zen" {
		t.Errorf("expected provider 'opencode-zen', got '%s'", pickle.Provider)
	}
	if pickle.SupportsVision == nil || *pickle.SupportsVision {
		t.Error("expected opencode-zen/big-pickle supports_vision to be false")
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

func isNil(v interface{}) bool {
	return v == nil
}

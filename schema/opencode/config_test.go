// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package opencode_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tforceaio/llm-db/schema/opencode"
)

func TestLoadConfig(t *testing.T) {
	data, err := os.ReadFile("../../json/samples/opencode.json")
	if err != nil {
		t.Fatalf("failed to read sample file: %v", err)
	}

	cfg, err := opencode.LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Schema != "https://opencode.ai/config.json" {
		t.Errorf("expected schema 'https://opencode.ai/config.json', got '%s'", cfg.Schema)
	}

	if cfg.SmallModel != "opencode/deepseek-v4-flash-free" {
		t.Errorf("expected small_model 'opencode/deepseek-v4-flash-free', got '%s'", cfg.SmallModel)
	}

	if cfg.DefaultAgent != "plan" {
		t.Errorf("expected default_agent 'plan', got '%s'", cfg.DefaultAgent)
	}

	if len(cfg.Agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(cfg.Agents))
	}

	if cfg.Agents["build"].Color != "primary" {
		t.Errorf("expected build agent color 'primary', got '%s'", cfg.Agents["build"].Color)
	}

	if cfg.Agents["plan"].Color != "secondary" {
		t.Errorf("expected plan agent color 'secondary', got '%s'", cfg.Agents["plan"].Color)
	}

	if len(cfg.MCPs) != 1 {
		t.Fatalf("expected 1 MCP config, got %d", len(cfg.MCPs))
	}

	bifrost := cfg.MCPs["bifrost"]
	if bifrost.Type != "remote" {
		t.Errorf("expected mcp type 'remote', got '%s'", bifrost.Type)
	}

	if !bifrost.Enabled {
		t.Error("expected mcp enabled to be true")
	}

	if len(cfg.Providers) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(cfg.Providers))
	}

	ollama := cfg.Providers["ollama"]
	if ollama.Name != "Ollama" {
		t.Errorf("expected provider name 'Ollama', got '%s'", ollama.Name)
	}

	gemma := ollama.Models["gemma-4:31b-instruct"]
	if gemma.ID != "gemma-4:31b-instruct" {
		t.Errorf("expected model id 'gemma-4:31b-instruct', got '%s'", gemma.ID)
	}

	if gemma.Cost == nil || gemma.Cost.Input != 0.140 {
		t.Errorf("expected cost input 0.140, got %v", gemma.Cost)
	}

	if gemma.Limit == nil || gemma.Limit.Context != 98304 {
		t.Errorf("expected limit context 98304, got %v", gemma.Limit)
	}

	if gemma.Modalities == nil || len(gemma.Modalities.Input) != 2 {
		t.Error("expected modalities input to have 2 elements")
	}

	if gemma.SupportReasoning == nil || !*gemma.SupportReasoning {
		t.Error("expected reasoning to be true")
	}
}

func TestRoundTrip(t *testing.T) {
	data, err := os.ReadFile("../../json/samples/opencode.json")
	if err != nil {
		t.Fatalf("failed to read sample file: %v", err)
	}

	cfg, err := opencode.LoadConfig(data)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	marshaled, err := json.MarshalIndent(cfg, "", "  ")
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

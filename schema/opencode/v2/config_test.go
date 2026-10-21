// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package v2_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tforceaio/llm-db/schema/opencode/v2"
)

const sampleConfig = `{
  "$schema": "https://opencode.ai/config.json",
  "model": "ollama/gemma-4:31b-instruct",
  "default_agent": "build",
  "agents": {
    "review": {
      "mode": "subagent",
      "color": "#00AA00",
      "steps": 3,
      "permissions": [{"action": "edit", "resource": "*", "effect": "allow"}]
    }
  },
  "mcp": {
    "timeout": {"startup": 5000},
    "servers": {
      "worker": {
        "type": "local",
        "command": ["bun", "x", "worker-mcp"],
        "cwd": "/repo",
        "environment": {"WORKER_TOKEN": "tok-00000000"},
        "protocol": "auto"
      },
      "bifrost": {
        "type": "remote",
        "url": "http://ollamago.tforce.xyz/mcp",
        "headers": {"x-bf-vk": "sk-bf-00000000"},
        "oauth": false,
        "codemode": false,
        "timeout": {"catalog": 2000, "execution": 3000}
      }
    }
  },
  "providers": {
    "ollama": {
      "name": "Ollama",
      "env": ["OLLAMA_API_KEY"],
      "package": "@ai-sdk/openai-compatible",
      "compatibility": {"reasoningField": "reasoning_content", "streamUsage": false},
      "settings": {"baseURL": "http://ollamago.tforce.xyz/v1", "apiKey": "sk-bf-00000000", "timeout": false},
      "headers": {"x-bf-vk": "sk-bf-00000000"},
      "body": {"stream_options": {"include_usage": true}},
      "models": {
        "gemma-4:31b-instruct": {
          "modelID": "gemma-4:31b-instruct",
          "name": "Gemma 4 31B Instruct",
          "releaseDate": "2026-09",
          "capabilities": {
            "tools": true,
            "input": ["text", "image"],
            "output": ["text"],
            "temperature": true,
            "reasoning": {"supported": true, "toggle": true, "budget": {"min": 256, "max": 8192}}
          },
          "cost": {"tier": {"type": "context", "size": 200000}, "input": 0.0000001, "output": 0.0000002, "cache": {"read": 0.00000005, "write": 0.0000004}},
          "limit": {"context": 98304, "input": 98304, "output": 16384},
          "disabled": false
        }
      }
    },
    "anthropic": {
      "compatibility": {"reasoningField": "reasoning", "requireReasoning": true},
      "models": {
        "claude-opus-4.7": {
          "compatibility": {"supportsThinkingBlockBinding": true, "supportsEffortUpdates": true},
          "variants": [{"id": "think", "headers": {"x-mode": "think"}}],
          "cost": [
            {"tier": {"type": "context", "size": 128000}, "input": 1.0, "output": 2.0},
            {"tier": {"type": "context", "size": 512000}, "input": 3.0, "output": 6.0}
          ]
        }
      }
    }
  }
}`

func TestLoadConfig(t *testing.T) {
	cfg, err := v2.LoadConfig([]byte(sampleConfig))
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Model == nil || !cfg.Model.ShortForm || cfg.Model.ProviderID != "ollama" || cfg.Model.Model != "gemma-4:31b-instruct" {
		t.Errorf("expected short-form model ollama/gemma-4:31b-instruct, got %+v", cfg.Model)
	}

	if cfg.DefaultAgent != "build" {
		t.Errorf("expected default_agent 'build', got '%s'", cfg.DefaultAgent)
	}

	review := cfg.Agents["review"]
	if review.Mode != "subagent" || review.Color != "#00AA00" {
		t.Errorf("expected review agent mode 'subagent' colored '#00AA00', got %+v", review)
	}

	if len(review.Permissions) != 1 || review.Permissions[0].Effect != "allow" {
		t.Errorf("expected one edit permission with effect 'allow', got %+v", review.Permissions)
	}

	if cfg.MCP == nil || cfg.MCP.Timeout == nil || cfg.MCP.Timeout.Startup == nil || *cfg.MCP.Timeout.Startup != 5000 {
		t.Errorf("expected mcp startup timeout 5000, got %+v", cfg.MCP)
	}

	worker := cfg.MCP.Servers["worker"]
	if worker.Type != "local" || len(worker.Command) != 3 || worker.Environment["WORKER_TOKEN"] != "tok-00000000" {
		t.Errorf("expected local worker server with command and token, got %+v", worker)
	}

	bifrost := cfg.MCP.Servers["bifrost"]
	if bifrost.Type != "remote" || bifrost.URL != "http://ollamago.tforce.xyz/mcp" {
		t.Errorf("expected remote bifrost server at ollamago.tforce.xyz/mcp, got %+v", bifrost)
	}

	if bifrost.OAuth == nil || !bifrost.OAuth.Disabled {
		t.Errorf("expected bifrost oauth explicitly disabled, got %+v", bifrost.OAuth)
	}

	ollama := cfg.Providers["ollama"]
	if ollama.Settings == nil || ollama.Settings["baseURL"] != "http://ollamago.tforce.xyz/v1" || ollama.Settings["timeout"] != false {
		t.Errorf("expected ollama settings with baseURL and disabled timeout, got %+v", ollama.Settings)
	}

	if ollama.Compatibility == nil || ollama.Compatibility.StreamUsage == nil || *ollama.Compatibility.StreamUsage {
		t.Errorf("expected ollama streamUsage disabled, got %+v", ollama.Compatibility)
	}

	gemma := ollama.Models["gemma-4:31b-instruct"]
	if gemma.Cost == nil || len(gemma.Cost) != 1 || gemma.Cost[0].Input == nil || *gemma.Cost[0].Input != 0.0000001 {
		t.Errorf("expected single gemma cost with input 0.0000001, got %+v", gemma.Cost)
	}

	if gemma.Cost[0].Tier == nil || gemma.Cost[0].Tier.Size != 200000 || gemma.Limit == nil || gemma.Limit.Input == nil || *gemma.Limit.Input != 98304 {
		t.Errorf("expected gemma cost tier of 200000 tokens and limit input 98304, got tier %+v limit %+v", gemma.Cost[0].Tier, gemma.Limit)
	}

	if gemma.Cost[0].Cache == nil || gemma.Cost[0].Cache.Read == nil || *gemma.Cost[0].Cache.Read != 0.00000005 {
		t.Errorf("expected cache read cost 0.00000005, got %+v", gemma.Cost[0].Cache)
	}

	if gemma.Capabilities == nil || gemma.Capabilities.Reasoning == nil || !gemma.Capabilities.Reasoning.Supported {
		t.Errorf("expected gemma reasoning supported, got %+v", gemma.Capabilities)
	}

	opus := cfg.Providers["anthropic"].Models["claude-opus-4.7"]
	if opus.Cost == nil || len(opus.Cost) != 2 {
		t.Fatalf("expected claude-opus costs packed as two tiers, got %+v", opus.Cost)
	}

	if opus.Cost[0].Tier == nil || opus.Cost[0].Tier.Size != 128000 {
		t.Errorf("expected first claude-opus cost tier of 128000 tokens, got %+v", opus.Cost[0].Tier)
	}
}

func TestRoundTrip(t *testing.T) {
	cfg, err := v2.LoadConfig([]byte(sampleConfig))
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	marshaled, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var original map[string]interface{}
	if err := json.Unmarshal([]byte(sampleConfig), &original); err != nil {
		t.Fatalf("failed to unmarshal original: %v", err)
	}

	var roundtripped map[string]interface{}
	if err := json.Unmarshal(marshaled, &roundtripped); err != nil {
		t.Fatalf("failed to unmarshal roundtripped: %v", err)
	}

	assertJSONEqual(t, original, roundtripped, "original", "roundtripped")
}

func TestModelSelectionReference(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    v2.ModelSelection
		wantErr bool
	}{
		{
			name:  "resolves the provider/model reference form",
			input: `"ollama/gemma-4:31b-instruct"`,
			want:  v2.ModelSelection{ProviderID: "ollama", Model: "gemma-4:31b-instruct", ShortForm: true},
		},
		{
			name:  "resolves the variant reference form",
			input: `"anthropic/claude-opus-4.7#think"`,
			want:  v2.ModelSelection{ProviderID: "anthropic", Model: "claude-opus-4.7", Variant: "think", ShortForm: true},
		},
		{
			name:  "keeps the explicit identifiers object form",
			input: `{"providerID": "ollama", "model": "gemma-4:31b-instruct"}`,
			want:  v2.ModelSelection{ProviderID: "ollama", Model: "gemma-4:31b-instruct", ShortForm: false},
		},
		{
			name:    "rejects a reference without a provider/model pair",
			input:   `"ollama"`,
			wantErr: true,
		},
		{
			name:    "rejects a provider without a model",
			input:   `{"providerID": "ollama"}`,
			wantErr: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeSelection(t, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected decoding %s to fail, got %v", tt.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected decoding %s to succeed: %v", tt.input, err)
			}
			if *got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, *got)
			}
			encoded, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if normalizeJSON(t, string(encoded)) != normalizeJSON(t, tt.input) {
				t.Errorf("expected %s to re-encode as itself, got %s", tt.input, encoded)
			}
		})
	}
}

func TestMCPOAuthOptOut(t *testing.T) {
	var explicitFalse v2.MCPOAuth
	if err := json.Unmarshal([]byte("false"), &explicitFalse); err != nil {
		t.Fatalf("expected oauth false to decode: %v", err)
	}
	if !explicitFalse.Disabled || explicitFalse.Config != nil {
		t.Errorf("expected oauth false to opt out, got %+v", explicitFalse)
	}

	encoded, err := json.Marshal(explicitFalse)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if string(encoded) != "false" {
		t.Errorf("expected opted-out oauth to encode as false, got %s", encoded)
	}

	var literalTrue v2.MCPOAuth
	if err := json.Unmarshal([]byte("true"), &literalTrue); err == nil {
		t.Error("expected oauth true to be rejected")
	}

	var zero v2.MCPOAuth
	if _, err := json.Marshal(zero); err == nil {
		t.Error("expected zero oauth to fail encoding instead of writing null")
	}
}

func decodeSelection(t *testing.T, input string) (*v2.ModelSelection, error) {
	t.Helper()
	var got v2.ModelSelection
	if err := json.Unmarshal([]byte(input), &got); err != nil {
		return nil, err
	}
	return &got, nil
}

// normalizeJSON canonicalizes a JSON document so forms compare by content and
// key order does not matter.
func normalizeJSON(t *testing.T, input string) string {
	t.Helper()
	var normalized interface{}
	if err := json.Unmarshal([]byte(input), &normalized); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", input, err)
	}
	reencoded, err := json.Marshal(normalized)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	return string(reencoded)
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

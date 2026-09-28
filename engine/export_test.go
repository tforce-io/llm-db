// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package engine

import (
	"testing"

	"github.com/tforceaio/llm-db/common"
	"github.com/tforceaio/llm-db/schema/llmdb"
)

func TestBuildBifrostModelsMode(t *testing.T) {
	m := &ExportModule{}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"m-embed": {
				Capabilities: []string{llmdb.CapabilityEmbedding},
				Deployments:  map[string]llmdb.Deployment{"localai": {ID: "m-embed"}},
			},
			"m-imggen": {
				Capabilities: []string{llmdb.CapabilityImageGeneration},
				Deployments:  map[string]llmdb.Deployment{"localai": {ID: "m-imggen"}},
			},
			"m-imgedit": {
				Capabilities: []string{llmdb.CapabilityImageEdit},
				Deployments:  map[string]llmdb.Deployment{"localai": {ID: "m-imgedit"}},
			},
			"m-chat": {
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments:  map[string]llmdb.Deployment{"localai": {ID: "m-chat"}},
			},
			"m-multi": {
				Capabilities: []string{llmdb.CapabilityImageEdit, llmdb.CapabilityImageGeneration},
				Deployments:  map[string]llmdb.Deployment{"localai": {ID: "m-multi"}},
			},
		},
	}

	result := m.buildBifrostModels(models)

	cases := map[string]string{
		"m-embed":   "embedding",
		"m-imggen":  "image_generation",
		"m-imgedit": "image_edit",
		"m-chat":    "chat",
		"m-multi":   "image_generation",
	}

	for key, expected := range cases {
		bModel, ok := result[key]
		if !ok {
			t.Fatalf("expected model %q in result", key)
		}
		if bModel.Mode != expected {
			t.Errorf("%s: expected mode %q, got %q", key, expected, bModel.Mode)
		}
	}
}

func TestBuildBifrostModelsVariantKeys(t *testing.T) {
	m := &ExportModule{}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"glm-4.7-flash-30b-a3b": {
				Name:         "GLM 4.7 Flash 30B A3B",
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments: map[string]llmdb.Deployment{
					"opencode-go": {ID: "glm-4-7-flash"},
					"localai": {
						Variants: []*llmdb.Deployment{
							{ID: "glm-4.7-flash-30b-a3b-fp8"},
							{ID: "glm-4.7-flash-30b-a3b-q4"},
						},
					},
				},
			},
			"some-model": {
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments: map[string]llmdb.Deployment{
					"groq": {
						Variants: []*llmdb.Deployment{
							{ID: "some-model-free"},
							{ID: ""},
						},
					},
				},
			},
		},
	}

	result := m.buildBifrostModels(models)

	expected := map[string]string{
		"opencode-go/glm-4-7-flash": "opencode-go",
		"glm-4.7-flash-30b-a3b-fp8": "localai",
		"glm-4.7-flash-30b-a3b-q4":  "localai",
		"groq/some-model-free":      "groq",
	}
	for key, provider := range expected {
		bModel, ok := result[key]
		if !ok {
			t.Fatalf("expected model %q in result", key)
		}
		if bModel.Provider != provider {
			t.Errorf("%s: expected provider %q, got %q", key, provider, bModel.Provider)
		}
	}
	if len(result) != len(expected) {
		t.Errorf("expected %d exported models, got %d", len(expected), len(result))
	}

	for key := range result {
		if key == "" || key[len(key)-1] == '/' {
			t.Errorf("unexpected malformed export key %q", key)
		}
	}
}

func TestBuildBifrostModelsCapabilityOverride(t *testing.T) {
	m := &ExportModule{}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"m-chat": {
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments: map[string]llmdb.Deployment{
					"localai": {
						Variants: []*llmdb.Deployment{
							{ID: "m-embed-v", Capabilities: []string{llmdb.CapabilityEmbedding}},
						},
					},
				},
			},
		},
	}

	result := m.buildBifrostModels(models)

	bModel, ok := result["m-embed-v"]
	if !ok {
		t.Fatalf("expected model %q in result", "m-embed-v")
	}
	if bModel.Mode != "embedding" {
		t.Errorf("expected mode %q from deployment capabilities, got %q", "embedding", bModel.Mode)
	}
	if bModel.SupportsReasoning != nil {
		t.Errorf("expected model-level reasoning flag overridden away, got %v", *bModel.SupportsReasoning)
	}
}

func TestBuildBifrostModelsCollision(t *testing.T) {
	m := &ExportModule{}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"m-a": {
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments:  map[string]llmdb.Deployment{"opencode-go": {ID: "same-id"}},
			},
			"m-b": {
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments:  map[string]llmdb.Deployment{"opencode-go": {ID: "same-id"}},
			},
		},
	}

	result := m.buildBifrostModels(models)

	if len(result) != 1 {
		t.Errorf("expected 1 exported model after collision drop, got %d", len(result))
	}
	if _, ok := result["opencode-go/same-id"]; !ok {
		t.Errorf("expected surviving export key opencode-go/same-id")
	}
}

func TestBuildOpenCodeConfigChatOnly(t *testing.T) {
	m := &ExportModule{}

	buildModel := func(capabilities []string) llmdb.Model {
		return llmdb.Model{
			Name:         "Model",
			Capabilities: capabilities,
			Deployments: map[string]llmdb.Deployment{
				"localai":    {ID: "id"},
				"openrouter": {ID: "id"},
			},
		}
	}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"m-chat":    buildModel([]string{llmdb.CapabilityReasoning}),
			"m-embed":   buildModel([]string{llmdb.CapabilityEmbedding}),
			"m-imggen":  buildModel([]string{llmdb.CapabilityImageGeneration}),
			"m-imgedit": buildModel([]string{llmdb.CapabilityImageEdit}),
		},
	}

	providers := &llmdb.Providers{Providers: map[string]llmdb.Provider{
		"localai": {Name: "LocalAI"},
		"bifrost": {Name: "Bifrost"},
	}}

	cfg := m.buildOpenCodeConfig(models, providers, "", "", "", "")

	expected := map[string][]string{
		"localai": {"id"},
		"bifrost": {"openrouter/id"},
	}
	excluded := []string{"m-embed", "m-imggen", "m-imgedit"}

	for providerKey, keys := range expected {
		provider, ok := cfg.Providers[providerKey]
		if !ok {
			t.Fatalf("expected provider %q in config", providerKey)
		}

		for _, key := range keys {
			if _, ok := provider.Models[key]; !ok {
				t.Errorf("%s: expected chat model %q to be included", providerKey, key)
			}
		}

		for _, name := range excluded {
			if _, ok := provider.Models[name]; ok {
				t.Errorf("%s: expected non-chat model %q to be excluded", providerKey, name)
			}
		}
	}
}

func TestBuildOpenCodeConfigVariantNaming(t *testing.T) {
	m := &ExportModule{}

	models := &llmdb.Models{
		Models: map[string]llmdb.Model{
			"glm-4.7-flash-30b-a3b": {
				Name:         "GLM 4.7 Flash 30B A3B",
				Capabilities: []string{llmdb.CapabilityReasoning},
				Deployments: map[string]llmdb.Deployment{
					"opencode-go": {ID: "glm-4-7-flash"},
					"localai": {
						Variants: []*llmdb.Deployment{
							{ID: "glm-fp8", Name: "fp8"},
							{ID: "glm-q4"},
						},
					},
				},
			},
		},
	}

	providers := &llmdb.Providers{Providers: map[string]llmdb.Provider{
		"bifrost":     {Name: "Bifrost"},
		"localai":     {Name: "LocalAI"},
		"opencode-go": {Name: "OpenCode Go"},
	}}

	cfg := m.buildOpenCodeConfig(models, providers, "", "", "", "")

	bifrostModel, ok := cfg.Providers["bifrost"].Models["opencode-go/glm-4-7-flash"]
	if !ok {
		t.Fatalf("expected bifrost model opencode-go/glm-4-7-flash")
	}
	if bifrostModel.Name != "GLM 4.7 Flash 30B A3B (OpenCode Go)" {
		t.Errorf("main deployment name: want %q, got %q", "GLM 4.7 Flash 30B A3B (OpenCode Go)", bifrostModel.Name)
	}

	fp8Model, ok := cfg.Providers["localai"].Models["glm-fp8"]
	if !ok {
		t.Fatalf("expected localai model glm-fp8")
	}
	if fp8Model.Name != "GLM 4.7 Flash 30B A3B [fp8]" {
		t.Errorf("variant deployment.Name: want %q, got %q", "GLM 4.7 Flash 30B A3B [fp8]", fp8Model.Name)
	}

	q4Model, ok := cfg.Providers["localai"].Models["glm-q4"]
	if !ok {
		t.Fatalf("expected localai model glm-q4")
	}
	if q4Model.Name != "GLM 4.7 Flash 30B A3B" {
		t.Errorf("variant without Name: want %q, got %q", "GLM 4.7 Flash 30B A3B", q4Model.Name)
	}
}

func TestResolveCostWithFactor(t *testing.T) {
	m := &ExportModule{}

	f15 := common.Float64(1.5)
	f25 := common.Float64(2.5)
	f0 := common.Float64(0)

	base := llmdb.ModelCost{
		TextInput:  common.Float64Ptr(0.070),
		TextOutput: common.Float64Ptr(0.400),
		CacheRead:  common.Float64Ptr(0.070),
	}

	cases := []struct {
		name       string
		override   *llmdb.ModelCost
		wantInput  *common.Float64
		wantOutput *common.Float64
		wantCache  *common.Float64
	}{
		{
			name:       "nil override keeps base",
			wantInput:  fp(0.070),
			wantOutput: fp(0.400),
			wantCache:  fp(0.070),
		},
		{
			name:       "override without factor replaces base",
			override:   &llmdb.ModelCost{TextInput: common.Float64Ptr(9.0)},
			wantInput:  fp(9.0),
			wantOutput: nil,
			wantCache:  nil,
		},
		{
			name:       "factor scales base",
			override:   &llmdb.ModelCost{Factor: &f15},
			wantInput:  fp(common.Float64(float64(*base.TextInput) * float64(f15))),
			wantOutput: fp(common.Float64(float64(*base.TextOutput) * float64(f15))),
			wantCache:  fp(common.Float64(float64(*base.CacheRead) * float64(f15))),
		},
		{
			name:       "factor wins over explicit fields",
			override:   &llmdb.ModelCost{Factor: &f25, TextInput: common.Float64Ptr(9.0)},
			wantInput:  fp(common.Float64(float64(*base.TextInput) * float64(f25))),
			wantOutput: fp(common.Float64(float64(*base.TextOutput) * float64(f25))),
			wantCache:  fp(common.Float64(float64(*base.CacheRead) * float64(f25))),
		},
		{
			name:       "zero factor zeroes present fields",
			override:   &llmdb.ModelCost{Factor: &f0},
			wantInput:  fp(0),
			wantOutput: fp(0),
			wantCache:  fp(0),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cost := m.resolveCost(base, c.override)
			assertCostField(t, "input", cost.TextInput, c.wantInput)
			assertCostField(t, "output", cost.TextOutput, c.wantOutput)
			assertCostField(t, "cache_read", cost.CacheRead, c.wantCache)
		})
	}
}

func TestResolveBifrostMode(t *testing.T) {
	cases := []struct {
		name         string
		capabilities []string
		expected     string
	}{
		{name: "empty", capabilities: nil, expected: "chat"},
		{name: "chat", capabilities: []string{llmdb.CapabilityReasoning}, expected: "chat"},
		{name: "embedding", capabilities: []string{llmdb.CapabilityEmbedding}, expected: "embedding"},
		{name: "image_generation", capabilities: []string{llmdb.CapabilityImageGeneration}, expected: "image_generation"},
		{name: "image_edit", capabilities: []string{llmdb.CapabilityImageEdit}, expected: "image_edit"},
		{name: "precedence image_generation over image_edit", capabilities: []string{llmdb.CapabilityImageEdit, llmdb.CapabilityImageGeneration}, expected: "image_generation"},
		{name: "precedence embedding over image", capabilities: []string{llmdb.CapabilityImageGeneration, llmdb.CapabilityEmbedding}, expected: "embedding"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := modelMode(c.capabilities); got != c.expected {
				t.Errorf("expected mode %q, got %q", c.expected, got)
			}
		})
	}
}

func fp(f common.Float64) *common.Float64 {
	return &f
}

func assertCostField(t *testing.T, label string, got, want *common.Float64) {
	t.Helper()

	switch {
	case want == nil && got != nil:
		t.Errorf("%s: want <nil>, got %v", label, float64(*got))
	case want != nil && got == nil:
		t.Errorf("%s: want %v, got <nil>", label, float64(*want))
	case want != nil && float64(*got) != float64(*want):
		t.Errorf("%s: want %v, got %v", label, float64(*want), float64(*got))
	}
}

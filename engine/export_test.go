// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package engine

import (
	"testing"

	"github.com/tforceaio/llm-db/schema/llmdb"
)

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

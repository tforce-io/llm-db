// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package llmdb_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/tforceaio/llm-db/schema/llmdb"
)

func TestLoadProviders(t *testing.T) {
	data, err := os.ReadFile("../../json/providers.json")
	if err != nil {
		t.Fatalf("failed to read providers file: %v", err)
	}

	providers, err := llmdb.LoadProviders(data)
	if err != nil {
		t.Fatalf("LoadProviders failed: %v", err)
	}

	if len(providers.Providers) == 0 {
		t.Fatal("expected at least one provider entry")
	}

	for id, p := range providers.Providers {
		t.Run(id, func(t *testing.T) {
			assertProviderValid(t, id, p)
		})
	}
}

func assertProviderValid(t *testing.T, id string, p llmdb.Provider) {
	t.Helper()

	if p.Name == "" {
		t.Errorf("%s: name is required", id)
	}
}

func TestProvidersRoundTrip(t *testing.T) {
	data, err := os.ReadFile("../../json/providers.json")
	if err != nil {
		t.Fatalf("failed to read providers file: %v", err)
	}

	providers, err := llmdb.LoadProviders(data)
	if err != nil {
		t.Fatalf("LoadProviders failed: %v", err)
	}

	marshaled, err := json.MarshalIndent(providers, "", "  ")
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

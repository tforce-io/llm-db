// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package llmdb_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tforceaio/llm-db/schema/llmdb"
)

func TestLoadModels(t *testing.T) {
	data, err := os.ReadFile("../../json/models.json")
	if err != nil {
		t.Fatalf("failed to read models file: %v", err)
	}

	models, err := llmdb.LoadModels(data)
	if err != nil {
		t.Fatalf("LoadModels failed: %v", err)
	}

	if len(models.Models) == 0 {
		t.Fatal("expected at least one model entry")
	}

	for id, m := range models.Models {
		t.Run(id, func(t *testing.T) {
			assertModelValid(t, id, m)
		})
	}
}

func assertModelValid(t *testing.T, id string, m llmdb.Model) {
	t.Helper()

	if m.Name == "" {
		t.Errorf("%s: name is required", id)
	}

	if m.Home == "" {
		t.Errorf("%s: home is required", id)
	}

	if len(m.Capabilities) == 0 {
		t.Errorf("%s: capabilities must not be empty", id)
	} else {
		for _, cap := range m.Capabilities {
			if !llmdb.ValidCapabilities[cap] {
				t.Errorf("%s: unknown capability %q", id, cap)
			}
		}
	}

	if m.Cost.Input < 0 {
		t.Errorf("%s: cost.input must not be negative", id)
	}

	if m.Cost.Output < 0 {
		t.Errorf("%s: cost.output must not be negative", id)
	}

	if m.Limit.Context <= 0 {
		t.Errorf("%s: limit.context must be positive", id)
	}

	if m.Limit.Output <= 0 {
		t.Errorf("%s: limit.output must be positive", id)
	}

	if len(m.Modalities.Input) == 0 {
		t.Errorf("%s: modalities.input must not be empty", id)
	} else {
		for _, mod := range m.Modalities.Input {
			if !llmdb.ValidModalities[mod] {
				t.Errorf("%s: unknown input modality %q", id, mod)
			}
		}
	}

	if len(m.Modalities.Output) == 0 {
		t.Errorf("%s: modalities.output must not be empty", id)
	} else {
		for _, mod := range m.Modalities.Output {
			if !llmdb.ValidModalities[mod] {
				t.Errorf("%s: unknown output modality %q", id, mod)
			}
		}
	}

	if len(m.Deployments) == 0 {
		t.Errorf("%s: deployment must not be empty", id)
	} else {
		for providerName, dep := range m.Deployments {
			if dep.ID == "" {
				t.Errorf("%s.deployment.%s: id is required", id, providerName)
			}
			if dep.Limit != nil {
				if dep.Limit.Context <= 0 {
					t.Errorf("%s.deployment.%s: limit.context must be positive", id, providerName)
				}
				if dep.Limit.Output <= 0 {
					t.Errorf("%s.deployment.%s: limit.output must be positive", id, providerName)
				}
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	data, err := os.ReadFile("../../json/models.json")
	if err != nil {
		t.Fatalf("failed to read models file: %v", err)
	}

	models, err := llmdb.LoadModels(data)
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

// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package llmdb

import "encoding/json"

type Providers struct {
	Schema    string             `json:"$schema,omitempty"`
	Providers map[string]Provider
}

func (p *Providers) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if schemaRaw, ok := raw["$schema"]; ok {
		if err := json.Unmarshal(schemaRaw, &p.Schema); err != nil {
			return err
		}
		delete(raw, "$schema")
	}

	p.Providers = make(map[string]Provider)
	for key, value := range raw {
		var provider Provider
		if err := json.Unmarshal(value, &provider); err != nil {
			return err
		}
		p.Providers[key] = provider
	}

	return nil
}

func (p Providers) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	if p.Schema != "" {
		result["$schema"] = p.Schema
	}
	for key, provider := range p.Providers {
		result[key] = provider
	}
	return json.Marshal(result)
}

type Provider struct {
	Name string `json:"name"`
	NPM  string `json:"npm,omitempty"`
	URI  string `json:"uri,omitempty"`
}

func LoadProviders(data []byte) (*Providers, error) {
	var providers Providers
	if err := json.Unmarshal(data, &providers); err != nil {
		return nil, err
	}
	return &providers, nil
}

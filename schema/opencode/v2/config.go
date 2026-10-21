// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

// Package v2 mirrors the OpenCode v2 configuration format (opencode.json) as published at https://opencode.ai/config.json.
package v2

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/tforceaio/llm-db/common"
)

// RootConfig is the root of an OpenCode v2 opencode.json document, narrowed to the sections this repository's exporter writes.
type RootConfig struct {
	Schema       string                    `json:"$schema,omitempty"`
	Model        *ModelSelection           `json:"model,omitempty"`
	DefaultAgent string                    `json:"default_agent,omitempty"`
	Agents       map[string]AgentConfig    `json:"agents,omitempty"`
	MCP          *MCPConfig                `json:"mcp,omitempty"`
	Providers    map[string]ProviderConfig `json:"providers,omitempty"`
}

// ModelSelection names a default model either as a provider/model reference. The loaded form decides the encoded form.
type ModelSelection struct {
	ProviderID string
	Model      string
	Variant    string
	ShortForm  bool
}

var (
	shortSelectionPattern = regexp.MustCompile(`^([^/#]+)/([^#]+)(?:#([^#]+))?$`)
	providerIDPattern     = regexp.MustCompile(`^[^/#]+$`)
	modelIDPattern        = regexp.MustCompile(`^[^#]+$`)
)

// Decode the provider/model reference form and the explicit identifiers object form.
func (s *ModelSelection) UnmarshalJSON(data []byte) error {
	var reference string
	if err := json.Unmarshal(data, &reference); err == nil {
		groups := shortSelectionPattern.FindStringSubmatch(reference)
		if groups == nil {
			return fmt.Errorf("parsing model selection %q: expected provider/model reference", reference)
		}
		s.ProviderID, s.Model, s.Variant, s.ShortForm = groups[1], groups[2], groups[3], true
		return nil
	}
	var explicit struct {
		ProviderID string `json:"providerID"`
		Model      string `json:"model"`
		Variant    string `json:"variant"`
	}
	if err := json.Unmarshal(data, &explicit); err != nil {
		return fmt.Errorf("parsing model selection: %w", err)
	}
	if !providerIDPattern.MatchString(explicit.ProviderID) || !modelIDPattern.MatchString(explicit.Model) {
		return fmt.Errorf("parsing model selection: providerID and model must be non-empty; providerID excludes '/' and '#'; model excludes '#'")
	}
	if explicit.Variant != "" && !modelIDPattern.MatchString(explicit.Variant) {
		return fmt.Errorf("parsing model selection: variant must not contain '#': %q", explicit.Variant)
	}
	s.ProviderID, s.Model, s.Variant, s.ShortForm = explicit.ProviderID, explicit.Model, explicit.Variant, false
	return nil
}

// Encode the selection back in the form it was decoded from.
func (s ModelSelection) MarshalJSON() ([]byte, error) {
	if s.ProviderID == "" || s.Model == "" {
		return nil, fmt.Errorf("unable to encode model selection: providerID and model are required")
	}
	if s.ShortForm {
		reference := s.ProviderID + "/" + s.Model
		if s.Variant != "" {
			return json.Marshal(reference + "#" + s.Variant)
		}
		return json.Marshal(reference)
	}
	return json.Marshal(struct {
		ProviderID string `json:"providerID"`
		Model      string `json:"model"`
		Variant    string `json:"variant,omitempty"`
	}{s.ProviderID, s.Model, s.Variant})
}

// AgentConfig is a named built-in agent override or custom agent definition.
type AgentConfig struct {
	Model       *ModelSelection  `json:"model,omitempty"`
	Request     *RequestOverlays `json:"request,omitempty"`
	System      string           `json:"system,omitempty"`
	Description string           `json:"description,omitempty"`
	Mode        string           `json:"mode,omitempty"`
	Hidden      *bool            `json:"hidden,omitempty"`
	Color       string           `json:"color,omitempty"`
	Steps       *int             `json:"steps,omitempty"`
	Disabled    *bool            `json:"disabled,omitempty"`
	Permissions []PermissionRule `json:"permissions,omitempty"`
}

// PermissionRule is one ordered allow/deny/ask rule applied to agent tool use.
type PermissionRule struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
	Effect   string `json:"effect"`
}

// RequestOverlays carries the request-level headers and body overrides.
type RequestOverlays struct {
	Headers map[string]string `json:"headers,omitempty"`
	Body    map[string]any    `json:"body,omitempty"`
}

// MCPConfig configures MCP servers, replacing the v1 per-server record with the named servers record.
type MCPConfig struct {
	Timeout *TimeoutConfig             `json:"timeout,omitempty"`
	Servers map[string]MCPServerConfig `json:"servers,omitempty"`
}

// TimeoutConfig bounds MCP operations in milliseconds.
type TimeoutConfig struct {
	Startup   *int `json:"startup,omitempty"`
	Catalog   *int `json:"catalog,omitempty"`
	Execution *int `json:"execution,omitempty"`
}

// MCPServerConfig holds either a local or remote server definition; the remaining fields apply per the set type.
type MCPServerConfig struct {
	Type        string            `json:"type"`
	Command     []string          `json:"command,omitempty"`
	Cwd         string            `json:"cwd,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	OAuth       *MCPOAuth         `json:"oauth,omitempty"`
	Disabled    *bool             `json:"disabled,omitempty"`
	CodeMode    *bool             `json:"codemode,omitempty"`
	Timeout     *TimeoutConfig    `json:"timeout,omitempty"`
	Protocol    string            `json:"protocol,omitempty"`
}

// MCPOAuth holds the server's OAuth configuration or an explicit opt-out written as false.
// The zero value cannot encode; name Config or set Disabled.
type MCPOAuth struct {
	Config   *MCPOAuthConfig
	Disabled bool
}

// Encode the explicit opt-out as false, or the OAuth configuration object.
func (o MCPOAuth) MarshalJSON() ([]byte, error) {
	if o.Disabled {
		return json.Marshal(false)
	}
	if o.Config == nil {
		return nil, fmt.Errorf("unable to encode mcp oauth: oauth config or explicit opt-out required")
	}
	return json.Marshal(o.Config)
}

// Decode the explicit opt-out literal (false) or an OAuth provider config.
func (o *MCPOAuth) UnmarshalJSON(data []byte) error {
	// A literal is present only when the config opts out; the schema allows
	// just the false literal, so true fails decoding.
	var literal bool
	if err := json.Unmarshal(data, &literal); err == nil {
		if literal {
			return fmt.Errorf("parsing mcp oauth: literal true is not a valid oauth configuration")
		}
		o.Config, o.Disabled = nil, true
		return nil
	}
	var cfg MCPOAuthConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parsing mcp oauth: %w", err)
	}
	o.Config, o.Disabled = &cfg, false
	return nil
}

// MCPOAuthConfig points an MCP server at an OAuth provider; every field is optional.
type MCPOAuthConfig struct {
	ClientID              string `json:"client_id,omitempty"`
	ClientSecret          string `json:"client_secret,omitempty"`
	Scope                 string `json:"scope,omitempty"`
	CallbackPort          *int   `json:"callback_port,omitempty"`
	RedirectURI           string `json:"redirect_uri,omitempty"`
	AuthServerMetadataURL string `json:"auth_server_metadata_url,omitempty"`
}

// ProviderConfig registers a provider, its package, and its model overrides.
type ProviderConfig struct {
	Canonical     string                 `json:"canonical,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Environment   []string               `json:"env,omitempty"`
	Package       string                 `json:"package,omitempty"`
	Compatibility *ProviderCompatibility `json:"compatibility,omitempty"`
	Settings      map[string]any         `json:"settings,omitempty"`
	Headers       map[string]string      `json:"headers,omitempty"`
	Body          map[string]any         `json:"body,omitempty"`
	Models        map[string]ModelConfig `json:"models,omitempty"`
}

// ProviderCompatibility carries per-provider transport quirks inherited by its models; the model-level fields vary, so the provider level omits them.
type ProviderCompatibility struct {
	ReasoningField            string `json:"reasoningField,omitempty"`
	RequireReasoning          *bool  `json:"requireReasoning,omitempty"`
	MaxTokensField            string `json:"maxTokensField,omitempty"`
	RequireFinishReason       *bool  `json:"requireFinishReason,omitempty"`
	RequireAssistantAfterTool *bool  `json:"requireAssistantAfterTool,omitempty"`
	SupportsPromptCacheKey    *bool  `json:"supportsPromptCacheKey,omitempty"`
	StreamUsage               *bool  `json:"streamUsage,omitempty"`
	CacheFormat               string `json:"cacheFormat,omitempty"`
}

// ModelConfig extends a provider's model; unset fields fall back to the provider defaults.
type ModelConfig struct {
	ModelID       string              `json:"modelID,omitempty"`
	Family        string              `json:"family,omitempty"`
	Name          string              `json:"name,omitempty"`
	ReleaseDate   string              `json:"releaseDate,omitempty"`
	Compatibility *ModelCompatibility `json:"compatibility,omitempty"`
	Package       string              `json:"package,omitempty"`
	Settings      map[string]any      `json:"settings,omitempty"`
	Headers       map[string]string   `json:"headers,omitempty"`
	Body          map[string]any      `json:"body,omitempty"`
	Capabilities  *Capabilities       `json:"capabilities,omitempty"`
	Variants      []ModelVariant      `json:"variants,omitempty"`
	Cost          ModelCosts          `json:"cost,omitempty"`
	Disabled      *bool               `json:"disabled,omitempty"`
	Limit         *ModelLimit         `json:"limit,omitempty"`
}

// ModelCompatibility carries per-model transport quirks affecting request encoding, cache behaviour, and thinking-block handling.
type ModelCompatibility struct {
	ReasoningField               string `json:"reasoningField,omitempty"`
	RequireReasoning             *bool  `json:"requireReasoning,omitempty"`
	MaxTokensField               string `json:"maxTokensField,omitempty"`
	RequireFinishReason          *bool  `json:"requireFinishReason,omitempty"`
	RequireAssistantAfterTool    *bool  `json:"requireAssistantAfterTool,omitempty"`
	SupportsPromptCacheKey       *bool  `json:"supportsPromptCacheKey,omitempty"`
	StreamUsage                  *bool  `json:"streamUsage,omitempty"`
	CacheFormat                  string `json:"cacheFormat,omitempty"`
	SupportsThinkingBlockBinding *bool  `json:"supportsThinkingBlockBinding,omitempty"`
	SupportsEffortUpdates        *bool  `json:"supportsEffortUpdates,omitempty"`
}

// Capabilities states what a model accepts; unset fields fall back to the provider defaults, then to the built-in defaults (tools with text and image input, text output).
type Capabilities struct {
	Tools       *bool      `json:"tools,omitempty"`
	Input       []string   `json:"input,omitempty"`
	Output      []string   `json:"output,omitempty"`
	Temperature *bool      `json:"temperature,omitempty"`
	Reasoning   *Reasoning `json:"reasoning,omitempty"`
}

// Reasoning is the models.dev v2 capabilities.reasoning record; Supported must always be written, while the extras only apply when reasoning is supported.
type Reasoning struct {
	Supported bool             `json:"supported"`
	Toggle    *bool            `json:"toggle,omitempty"`
	Effort    []string         `json:"effort,omitempty"`
	Budget    *ReasoningBudget `json:"budget,omitempty"`
}

// ReasoningBudget bounds the reasoning token budget.
type ReasoningBudget struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// ModelVariant scales a base model behind a variant ID with id-scoped overlays.
type ModelVariant struct {
	ID       string            `json:"id"`
	Settings map[string]any    `json:"settings,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     map[string]any    `json:"body,omitempty"`
}

// ModelCosts holds one cost table or a pack of per-context-tier costs; a single entry encodes as the object form, several as the array form.
type ModelCosts []ModelCost

// Decode either the single-object cost form or the per-context-tier array form.
func (c *ModelCosts) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		if err := json.Unmarshal(data, (*[]ModelCost)(c)); err != nil {
			return fmt.Errorf("parsing model costs: %w", err)
		}
		return nil
	}
	var single ModelCost
	if err := json.Unmarshal(data, &single); err != nil {
		return fmt.Errorf("parsing model costs: %w", err)
	}
	*c = ModelCosts{single}
	return nil
}

// Encode a single cost as an object and several costs as an array.
func (c ModelCosts) MarshalJSON() ([]byte, error) {
	if len(c) == 1 {
		return json.Marshal(c[0])
	}
	return json.Marshal([]ModelCost(c))
}

// ModelCost prices a model in USD per million tokens; input and output are required, tier and cache counts are optional.
type ModelCost struct {
	Tier   *ModelCostTier  `json:"tier,omitempty"`
	Input  *common.Float64 `json:"input"`
	Output *common.Float64 `json:"output"`
	Cache  *ModelCostCache `json:"cache,omitempty"`
}

// ModelCostTier prices an explicit context tier; Type is always "context" and Size counts the tokens the tier applies to.
type ModelCostTier struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

// ModelCostCache prices cached token reads and writes.
type ModelCostCache struct {
	Read  *common.Float64 `json:"read,omitempty"`
	Write *common.Float64 `json:"write,omitempty"`
}

// ModelLimit caps context, prompt input tokens, and output tokens.
type ModelLimit struct {
	Context *int `json:"context,omitempty"`
	Input   *int `json:"input,omitempty"`
	Output  *int `json:"output,omitempty"`
}

// Parse an OpenCode v2 configuration document; unknown keys are ignored, matching the v1 loader.
func LoadConfig(data []byte) (*RootConfig, error) {
	var cfg RootConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Package llm implements the embedded Claude assistant: an Anthropic Messages API
// client, a tool registry that wraps stash's data layer, and the agent loop that
// drives a conversation. See docs/llm/DESIGN.md.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	anthropicEndpoint = "https://api.anthropic.com/v1/messages"
	anthropicVersion  = "2023-06-01"
	defaultMaxTokens  = 4096
)

// Client is a minimal Anthropic Messages API client. It is intentionally
// dependency-free (net/http + encoding/json) so the fork stays lean and easy to
// merge with upstream.
type Client struct {
	apiKey    string
	model     string
	maxTokens int
	baseURL   string
	http      *http.Client
}

// NewClient builds a client for the given key/model. A fresh client is created per
// request so changes to the configured key/model take effect without a restart.
func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey:    apiKey,
		model:     model,
		maxTokens: defaultMaxTokens,
		baseURL:   anthropicEndpoint,
		http:      &http.Client{Timeout: 120 * time.Second},
	}
}

// Message is a single conversation message. Content is always the block form so we
// can carry tool_use / tool_result blocks alongside text.
type Message struct {
	Role    string         `json:"role"` // "user" | "assistant"
	Content []ContentBlock `json:"content"`
}

// ContentBlock is a tagged union over Anthropic content block types: "text",
// "tool_use", and "tool_result". Unused fields are omitted on marshal.
type ContentBlock struct {
	Type string `json:"type"`

	// type == "text"
	Text string `json:"text,omitempty"`

	// type == "tool_use"
	ID    string          `json:"id,omitempty"`
	Name  string          `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`

	// type == "tool_result"
	ToolUseID string `json:"tool_use_id,omitempty"`
	Content   string `json:"content,omitempty"`
	IsError   bool   `json:"is_error,omitempty"`
}

// ToolDef is a tool advertised to the model.
type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// Usage mirrors the token accounting returned by the API.
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system,omitempty"`
	Messages  []Message `json:"messages"`
	Tools     []ToolDef `json:"tools,omitempty"`
}

// Response is the decoded Messages API response.
type Response struct {
	ID         string         `json:"id"`
	Model      string         `json:"model"`
	Role       string         `json:"role"`
	StopReason string         `json:"stop_reason"`
	Content    []ContentBlock `json:"content"`
	Usage      Usage          `json:"usage"`
}

type apiError struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// CreateMessage performs a single (non-streaming) Messages API call.
func (c *Client) CreateMessage(ctx context.Context, system string, messages []Message, tools []ToolDef) (*Response, error) {
	reqBody := messagesRequest{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		System:    system,
		Messages:  messages,
		Tools:     tools,
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", anthropicVersion)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call anthropic: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var ae apiError
		if json.Unmarshal(body, &ae) == nil && ae.Error.Message != "" {
			return nil, fmt.Errorf("anthropic %d: %s", resp.StatusCode, ae.Error.Message)
		}
		return nil, fmt.Errorf("anthropic %d: %s", resp.StatusCode, string(body))
	}

	var out Response
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}

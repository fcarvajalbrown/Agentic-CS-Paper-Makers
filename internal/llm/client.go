package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.moonshot.cn/v1"
	maxRetries     = 3
)

type Client struct {
	apiKey  string
	baseURL string
	model   string
	seed    int
	http    *http.Client
	cache   *Cache
	tracker *CostTracker
}

type ClientOption func(*Client)

func WithBaseURL(url string) ClientOption {
	return func(c *Client) { c.baseURL = url }
}

func WithSeed(seed int) ClientOption {
	return func(c *Client) { c.seed = seed }
}

func WithCache(cache *Cache) ClientOption {
	return func(c *Client) { c.cache = cache }
}

func WithCostTracker(t *CostTracker) ClientOption {
	return func(c *Client) { c.tracker = t }
}

func NewClient(apiKey, model string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		model:   model,
		http:    &http.Client{Timeout: 120 * time.Second},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function ToolSpec `json:"function"`
}

type ToolSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

type ToolHandler func(name, arguments string) (string, error)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
	Seed     int       `json:"seed,omitempty"`
	Stream   bool      `json:"stream"`
}

type Response struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func (c *Client) Complete(stage string, messages []Message, tools []Tool, handler ToolHandler) (string, error) {
	for {
		resp, err := c.callWithRetry(messages, tools)
		if err != nil {
			return "", err
		}

		if c.tracker != nil {
			c.tracker.Record(stage, Usage{
				InputTokens:  resp.Usage.PromptTokens,
				OutputTokens: resp.Usage.CompletionTokens,
				TotalTokens:  resp.Usage.TotalTokens,
				CostUSD:      CostForModel(c.model, resp.Usage.PromptTokens, resp.Usage.CompletionTokens),
			})
		}

		choice := resp.Choices[0]

		if choice.FinishReason != "tool_calls" || len(choice.Message.ToolCalls) == 0 {
			return choice.Message.Content, nil
		}

		messages = append(messages, choice.Message)

		for _, tc := range choice.Message.ToolCalls {
			result, err := handler(tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				result = fmt.Sprintf("error: %v", err)
			}
			messages = append(messages, Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: tc.ID,
			})
		}
	}
}

func (c *Client) callWithRetry(messages []Message, tools []Tool) (*Response, error) {
	req := Request{
		Model:    c.model,
		Messages: messages,
		Tools:    tools,
		Seed:     c.seed,
	}

	var lastErr error
	for attempt := range maxRetries {
		resp, err := c.call(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		wait := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		time.Sleep(wait)
	}
	return nil, fmt.Errorf("after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) call(req Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", httpResp.StatusCode, string(raw))
	}

	var resp Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty choices in response")
	}
	return &resp, nil
}

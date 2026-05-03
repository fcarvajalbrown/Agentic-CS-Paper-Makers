package agents

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
)

// Runner executes a single agent turn against the LLM.
// It loads the agent's system prompt and prepends it to the provided message
// history before calling Complete. Callers maintain the message history for
// multi-turn agents (e.g. the Architect's Socratic loop).
type Runner struct {
	client *llm.Client
}

func NewRunner(client *llm.Client) *Runner {
	return &Runner{client: client}
}

// RunRequest carries everything needed for one agent call.
type RunRequest struct {
	// Agent identifies which profile to load as the system prompt.
	Agent AgentName

	// Stage is passed to the cost tracker (e.g. "research", "architect").
	Stage string

	// Messages is the conversation history not including the system message.
	// For single-turn agents, pass a slice with one user message.
	// For multi-turn agents (Architect), append previous turns before each call.
	Messages []llm.Message

	// Tools is the list of tools available to this agent.
	// Empty for all agents except the Research Agent.
	Tools []llm.Tool

	// Handler is called by the tool-use loop when the LLM invokes a tool.
	// May be nil when Tools is empty.
	Handler llm.ToolHandler
}

// Run loads the agent profile, prepends it as the system message, and calls
// the LLM. It returns the final text content of the response.
func (r *Runner) Run(req RunRequest) (string, error) {
	profile, err := LoadProfile(req.Agent)
	if err != nil {
		return "", fmt.Errorf("runner: %w", err)
	}

	msgs := make([]llm.Message, 0, 1+len(req.Messages))
	msgs = append(msgs, llm.Message{Role: "system", Content: profile})
	msgs = append(msgs, req.Messages...)

	handler := req.Handler
	if handler == nil {
		handler = func(name, arguments string) (string, error) {
			return "", fmt.Errorf("no tool handler registered for %q", name)
		}
	}

	return r.client.Complete(req.Stage, msgs, req.Tools, handler)
}

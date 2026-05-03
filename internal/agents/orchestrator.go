package agents

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
)

// Orchestrator dispatches agents sequentially or in parallel.
type Orchestrator struct {
	runner *Runner
}

func NewOrchestrator(runner *Runner) *Orchestrator {
	return &Orchestrator{runner: runner}
}

// RunSequential runs a single agent and returns its response text.
func (o *Orchestrator) RunSequential(req RunRequest) (string, error) {
	return o.runner.Run(req)
}

// CriticResult holds the outcome of one parallel critic run.
type CriticResult struct {
	Agent    AgentName
	Response string
	Err      error
}

// RunCritics dispatches Formalist, GameTheorist, and Skeptic concurrently.
// Each critic receives draftContent as its sole user message.
// All three results are returned even when some fail; the caller inspects
// CriticResult.Err per agent. The returned error is non-nil if any critic
// failed, combining all individual errors.
func (o *Orchestrator) RunCritics(draftContent string) ([3]CriticResult, error) {
	critics := [3]AgentName{AgentFormalist, AgentGameTheorist, AgentSkeptic}

	var (
		wg      sync.WaitGroup
		results [3]CriticResult
	)

	for i, agent := range critics {
		wg.Add(1)
		go func(idx int, a AgentName) {
			defer wg.Done()
			results[idx].Agent = a
			resp, err := o.runner.Run(RunRequest{
				Agent: a,
				Stage: string(a),
				Messages: []llm.Message{
					{Role: "user", Content: draftContent},
				},
			})
			results[idx].Response = resp
			results[idx].Err = err
		}(i, agent)
	}

	wg.Wait()

	var errs []error
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", r.Agent, r.Err))
		}
	}
	return results, errors.Join(errs...)
}

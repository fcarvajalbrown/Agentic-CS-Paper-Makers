package llm

import (
	"fmt"
	"strings"
)

const (
	costPer1MTokensKimiK25      = 0.15
	costPer1MTokensKimiK26      = 0.60
	costPer1MTokensMoonshotV132k = 0.08
)

type Usage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
	CostUSD      float64
}

func (u Usage) String() string {
	return fmt.Sprintf("tokens=%d cost=$%.4f", u.TotalTokens, u.CostUSD)
}

type CostTracker struct {
	stages map[string]Usage
}

func NewCostTracker() *CostTracker {
	return &CostTracker{stages: make(map[string]Usage)}
}

func (t *CostTracker) Record(stage string, u Usage) {
	prev := t.stages[stage]
	t.stages[stage] = Usage{
		InputTokens:  prev.InputTokens + u.InputTokens,
		OutputTokens: prev.OutputTokens + u.OutputTokens,
		TotalTokens:  prev.TotalTokens + u.TotalTokens,
		CostUSD:      prev.CostUSD + u.CostUSD,
	}
}

func (t *CostTracker) Total() Usage {
	var total Usage
	for _, u := range t.stages {
		total.InputTokens += u.InputTokens
		total.OutputTokens += u.OutputTokens
		total.TotalTokens += u.TotalTokens
		total.CostUSD += u.CostUSD
	}
	return total
}

func (t *CostTracker) Report() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-16s %10s %12s\n", "Stage", "Tokens", "Cost (USD)"))
	b.WriteString(strings.Repeat("-", 42) + "\n")
	for stage, u := range t.stages {
		b.WriteString(fmt.Sprintf("%-16s %10d %12s\n", stage, u.TotalTokens, fmt.Sprintf("$%.2f", u.CostUSD)))
	}
	b.WriteString(strings.Repeat("-", 42) + "\n")
	total := t.Total()
	b.WriteString(fmt.Sprintf("%-16s %10d %12s\n", "TOTAL", total.TotalTokens, fmt.Sprintf("$%.2f", total.CostUSD)))
	return b.String()
}

func CostForModel(model string, inputTokens, outputTokens int) float64 {
	rate := costPer1MTokensKimiK25
	switch model {
	case "kimi-k2.6":
		rate = costPer1MTokensKimiK26
	case "moonshot-v1-32k":
		rate = costPer1MTokensMoonshotV132k
	}
	return float64(inputTokens+outputTokens) / 1_000_000 * rate
}

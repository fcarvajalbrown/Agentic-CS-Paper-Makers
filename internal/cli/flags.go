package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/spf13/cobra"
)

// Global flag values — set by cobra before RunE executes.
var (
	flagModel      string
	flagBudget     float64
	flagMaxTokens  int
	flagCheap      bool
	flagProduction bool
	flagNoCache    bool
	flagDryRun     bool
	flagAgentModel []string // "agent:model" pairs, e.g. "architect:kimi-k2.6"
	flagSeed       string
)

// addGlobalFlags attaches flags that apply to all workflow commands.
func addGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&flagModel, "model", "", "override default model")
	cmd.PersistentFlags().Float64Var(&flagBudget, "budget", 0, "total spend cap in USD (e.g. 10.00)")
	cmd.PersistentFlags().IntVar(&flagMaxTokens, "max-tokens", 0, "hard token cap per LLM call")
	cmd.PersistentFlags().BoolVar(&flagCheap, "cheap", false, "use moonshot-v1-32k for all stages")
	cmd.PersistentFlags().BoolVar(&flagProduction, "production", false, "use kimi-k2.6 for all stages")
	cmd.PersistentFlags().BoolVar(&flagNoCache, "no-cache", false, "force fresh LLM call, skip cache")
	cmd.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "print prompt without sending to LLM")
	cmd.PersistentFlags().StringArrayVar(&flagAgentModel, "agent-model", nil, "per-agent model override (agent:model)")
}

// buildFlags converts the global flag values into a config.Flags struct.
func buildFlags() config.Flags {
	agentModels := make(map[string]string)
	for _, pair := range flagAgentModel {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			agentModels[parts[0]] = parts[1]
		}
	}
	return config.Flags{
		Model:      flagModel,
		Budget:     flagBudget,
		MaxTokens:  flagMaxTokens,
		Cheap:      flagCheap,
		Production: flagProduction,
		AgentModel: agentModels,
	}
}

// requireProject checks that the current directory contains a .paperflow/ directory
// and returns the absolute path of the current working directory.
func requireProject() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(cwd, ".paperflow")); os.IsNotExist(err) {
		return "", fmt.Errorf("not a paperflow project (no .paperflow/ found in %s)\nRun 'paperflow init <name>' to create one", cwd)
	}
	return cwd, nil
}

// newClient builds an LLM client for a given stage using the provided config and project root.
func newClient(cfg *config.Config, stage, cacheDir string) (*llm.Client, *llm.CostTracker, error) {
	if cfg.APIKey == "" {
		return nil, nil, fmt.Errorf("no API key — set %s or add apiKey to config", config.EnvAPIKey)
	}
	model := cfg.ModelForAgent(stage)
	tracker := llm.NewCostTracker()
	opts := []llm.ClientOption{
		llm.WithCostTracker(tracker),
		llm.WithSeed(42),
	}
	if !flagNoCache {
		opts = append(opts, llm.WithCache(llm.NewCache(cacheDir)))
	}
	return llm.NewClient(cfg.APIKey, model, opts...), tracker, nil
}

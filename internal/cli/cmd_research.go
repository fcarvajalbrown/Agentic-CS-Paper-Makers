package cli

import (
	"encoding/json"
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/tools"
	"github.com/spf13/cobra"
)

func newResearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "research",
		Short: "Run the Literature Scout agent (arXiv, Semantic Scholar, Zenodo)",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			cfg, err := config.Load(root, buildFlags())
			if err != nil {
				return err
			}

			seed := flagSeed
			if seed == "" {
				seed = cfg.Seed
			}
			if seed == "" {
				return fmt.Errorf("no seed found; use --seed or set it during 'paperflow init'")
			}

			client, _, err := newClient(cfg, "research", state.CacheDir(root))
			if err != nil {
				return err
			}

			runner := agents.NewRunner(client)

			researchTools := buildResearchTools()
			handler := buildResearchHandler()

			if flagDryRun {
				fmt.Printf("dry-run: would call research agent with seed: %q\n", seed)
				return nil
			}

			userMsg := fmt.Sprintf("Research seed: %s", seed)
			resp, err := runner.Run(agents.RunRequest{
				Agent:    agents.AgentResearch,
				Stage:    "research",
				Messages: []llm.Message{{Role: "user", Content: userMsg}},
				Tools:    researchTools,
				Handler:  handler,
			})
			if err != nil {
				return fmt.Errorf("research agent: %w", err)
			}

			ver := 1
			if err := artifacts.Write(state.ArtifactsDir(root), "research_context", ver, artifacts.SchemaResearchContext, []byte(resp)); err != nil {
				return fmt.Errorf("save research_context: %w", err)
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageResearch)
			cp.LatestArtifact = artifacts.ArtifactPath(state.ArtifactsDir(root), "research_context", ver)
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			fmt.Printf("Research complete. Artifact: research_context_v%d.json\n", ver)
			fmt.Println("Next step: paperflow architect")
			return nil
		},
	}
	cmd.Flags().StringVar(&flagSeed, "seed", "", "seed idea for the paper (overrides project config)")
	return cmd
}

func buildResearchTools() []llm.Tool {
	queryParam := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query":       map[string]any{"type": "string", "description": "search query string"},
			"max_results": map[string]any{"type": "integer", "description": "max papers to return", "default": 5},
		},
		"required": []string{"query"},
	}

	makeTool := func(name, desc string) llm.Tool {
		return llm.Tool{
			Type: "function",
			Function: llm.ToolSpec{
				Name:        name,
				Description: desc,
				Parameters:  queryParam,
			},
		}
	}

	return []llm.Tool{
		makeTool("search_arxiv", "Search arXiv preprint server for academic papers"),
		makeTool("search_semantic_scholar", "Search Semantic Scholar for academic papers"),
		makeTool("search_zenodo", "Search Zenodo for preprints and research datasets"),
	}
}

func buildResearchHandler() llm.ToolHandler {
	return func(name, arguments string) (string, error) {
		var args struct {
			Query      string `json:"query"`
			MaxResults int    `json:"max_results"`
		}
		if err := json.Unmarshal([]byte(arguments), &args); err != nil {
			return "", fmt.Errorf("parse tool args: %w", err)
		}
		if args.MaxResults == 0 {
			args.MaxResults = 5
		}

		switch name {
		case "search_arxiv":
			papers, err := tools.SearchArxiv(args.Query, args.MaxResults)
			if err != nil {
				return "", err
			}
			out, _ := json.Marshal(papers)
			return string(out), nil

		case "search_semantic_scholar":
			papers, err := tools.SearchSemanticScholar(args.Query, args.MaxResults)
			if err != nil {
				return "", err
			}
			out, _ := json.Marshal(papers)
			return string(out), nil

		case "search_zenodo":
			papers, err := tools.SearchZenodo(args.Query, args.MaxResults)
			if err != nil {
				return "", err
			}
			out, _ := json.Marshal(papers)
			return string(out), nil

		default:
			return "", fmt.Errorf("unknown tool: %s", name)
		}
	}
}

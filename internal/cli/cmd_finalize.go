package cli

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/export"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newFinalizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "finalize",
		Short: "Assemble paper.md, references.bib, and metadata.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			cfg, err := config.Load(root, buildFlags())
			if err != nil {
				return err
			}

			artDir := state.ArtifactsDir(root)
			draftVer, err := artifacts.LatestVersion(artDir, "draft")
			if err != nil {
				return fmt.Errorf("no draft found; run 'paperflow write' first")
			}
			draftData, err := artifacts.Read(artDir, "draft", draftVer)
			if err != nil {
				return err
			}

			client, _, err := newClient(cfg, "finalize", state.CacheDir(root))
			if err != nil {
				return err
			}

			if flagDryRun {
				fmt.Println("dry-run: would call finalize agent")
				return nil
			}

			runner := agents.NewRunner(client)
			resp, err := runner.Run(agents.RunRequest{
				Agent:    agents.AgentFinalize,
				Stage:    "finalize",
				Messages: []llm.Message{{Role: "user", Content: string(draftData)}},
			})
			if err != nil {
				return fmt.Errorf("finalize agent: %w", err)
			}

			var output export.FinalizeOutput
			if err := json.Unmarshal([]byte(resp), &output); err != nil {
				return fmt.Errorf("parse finalize output: %w", err)
			}

			finalDir := filepath.Join(artDir, "final")
			if err := export.WriteArtifacts(finalDir, &output); err != nil {
				return err
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageFinalize)
			cp.LatestArtifact = filepath.Join(finalDir, "paper.md")
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			fmt.Printf("Finalized:\n  %s/paper.md\n  %s/references.bib\n  %s/metadata.json\n", finalDir, finalDir, finalDir)
			fmt.Println("Next step: paperflow export --pdf  (optional)")
			return nil
		},
	}
}

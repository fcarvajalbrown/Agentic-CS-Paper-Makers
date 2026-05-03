package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newWriteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "write",
		Short: "Run the Lead Writer agent to generate the first draft",
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

			bpVer, err := artifacts.LatestVersion(artDir, "blueprint")
			if err != nil {
				return fmt.Errorf("no blueprint found; run 'paperflow architect' first")
			}

			rcVer, err := artifacts.LatestVersion(artDir, "research_context")
			if err != nil {
				return fmt.Errorf("no research_context found; run 'paperflow research' first")
			}

			bpData, err := artifacts.Read(artDir, "blueprint", bpVer)
			if err != nil {
				return err
			}
			rcData, err := artifacts.Read(artDir, "research_context", rcVer)
			if err != nil {
				return err
			}

			client, _, err := newClient(cfg, "writer", state.CacheDir(root))
			if err != nil {
				return err
			}

			if flagDryRun {
				fmt.Println("dry-run: would call writer agent")
				return nil
			}

			runner := agents.NewRunner(client)
			userMsg := fmt.Sprintf("Blueprint:\n%s\n\nResearch context:\n%s", string(bpData), string(rcData))

			resp, err := runner.Run(agents.RunRequest{
				Agent:    agents.AgentWriter,
				Stage:    "writer",
				Messages: []llm.Message{{Role: "user", Content: userMsg}},
			})
			if err != nil {
				return fmt.Errorf("writer agent: %w", err)
			}

			ver := 1
			if err := artifacts.Write(artDir, "draft", ver, artifacts.SchemaDraft, []byte(resp)); err != nil {
				return fmt.Errorf("save draft: %w", err)
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageWrite)
			cp.LatestArtifact = artifacts.ArtifactPath(artDir, "draft", ver)
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			fmt.Printf("Draft written: draft_v%d.json\n", ver)
			fmt.Println("Next step: paperflow review")
			return nil
		},
	}
}

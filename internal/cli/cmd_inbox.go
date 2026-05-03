package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/inbox"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newInboxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inbox",
		Short: "Interactive reviewer response loop",
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
			roundVer, err := artifacts.LatestVersion(artDir, "review_round")
			if err != nil {
				return fmt.Errorf("no review_round found; run 'paperflow review' first")
			}

			client, _, err := newClient(cfg, "writer", state.CacheDir(root))
			if err != nil {
				return err
			}

			runner := agents.NewRunner(client)
			ib := inbox.New(runner, artDir)

			if err := ib.Run(roundVer); err != nil {
				return err
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageInbox)
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			return nil
		},
	}
}

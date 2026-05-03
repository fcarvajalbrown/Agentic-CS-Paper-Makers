package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current stage and checkpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			cp, err := state.LoadCheckpoint(root)
			if err != nil {
				return err
			}

			if cp.Stage == "" {
				fmt.Println("Stage:           init (no commands run yet)")
			} else {
				fmt.Printf("Stage:           %s\n", cp.Stage)
			}
			if cp.LatestArtifact != "" {
				fmt.Printf("Latest artifact: %s\n", cp.LatestArtifact)
			}
			fmt.Printf("Total tokens:    %d\n", cp.TotalTokens)
			fmt.Printf("Total cost:      $%.4f\n", cp.TotalCost)
			return nil
		},
	}
}

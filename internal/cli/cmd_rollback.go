package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newRollbackCmd() *cobra.Command {
	var flagTo string

	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Revert checkpoint to an earlier stage",
		Long: `Revert the project checkpoint to an earlier stage.
The artifact files are preserved — rollback only resets the stage pointer.

Valid stages: init, research, architect, write, review, inbox, finalize`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagTo == "" {
				return fmt.Errorf("--to is required (e.g. --to=write)")
			}

			root, err := requireProject()
			if err != nil {
				return err
			}

			if _, err := state.ParseStage(flagTo); err != nil {
				return fmt.Errorf("invalid stage %q: %w", flagTo, err)
			}

			if err := state.RollbackTo(root, flagTo); err != nil {
				return fmt.Errorf("rollback: %w", err)
			}

			fmt.Printf("Checkpoint rolled back to stage: %s\n", flagTo)
			fmt.Println("Artifact files are unchanged. Re-run the stage command to regenerate.")
			return nil
		},
	}

	cmd.Flags().StringVar(&flagTo, "to", "", "target stage to roll back to")
	return cmd
}

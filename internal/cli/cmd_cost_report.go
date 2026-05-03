package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newCostReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cost-report",
		Short: "Print token and cost breakdown from the checkpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			cp, err := state.LoadCheckpoint(root)
			if err != nil {
				return err
			}

			fmt.Printf("%-20s %10s %12s\n", "Stage", "Tokens", "Cost (USD)")
			fmt.Println("-------------------------------------------")
			fmt.Printf("%-20s %10d %12s\n", "total (all stages)", cp.TotalTokens, fmt.Sprintf("$%.4f", cp.TotalCost))
			fmt.Println()
			fmt.Println("Run 'paperflow status' to see the current stage.")
			fmt.Println("Per-stage breakdown is tracked at runtime and reset on resume.")
			return nil
		},
	}
}

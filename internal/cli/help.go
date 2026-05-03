package cli

import (
	"github.com/spf13/cobra"
)

const usageHeader = `paperflow — multi-agent academic paper workflow

Workflow:
  paperflow init <name>   create a new paper project
  paperflow research      run Literature Scout (arXiv, Semantic Scholar, Zenodo)
  paperflow architect     Socratic framework builder (human gate)
  paperflow write         generate academic draft from approved blueprint
  paperflow review        dispatch 3 parallel critics (Formalist, Game Theorist, Skeptic)
  paperflow inbox         interactive reviewer response loop
  paperflow finalize      assemble paper.md, references.bib, metadata.json

Utilities:
  paperflow export        export to PDF (requires pandoc + wkhtmltopdf)
  paperflow status        show current stage and token spend
  paperflow resume        print next command to run after a crash
  paperflow rollback      revert to a previous artifact version
  paperflow cost-report   print spend breakdown by agent
  paperflow validate      validate any artifact JSON against its schema

Use 'paperflow <command> --help' for details on each command.
`

// NewRootCmd builds and returns the root cobra command with all subcommands attached.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "paperflow",
		Short:         "Multi-agent academic paper workflow",
		Long:          usageHeader,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	addGlobalFlags(root)

	root.AddCommand(
		newInitCmd(),
		newResearchCmd(),
		newArchitectCmd(),
		newWriteCmd(),
		newReviewCmd(),
		newInboxCmd(),
		newFinalizeCmd(),
		newExportCmd(),
		newStatusCmd(),
		newResumeCmd(),
		newRollbackCmd(),
		newCostReportCmd(),
		newValidateCmd(),
	)

	return root
}

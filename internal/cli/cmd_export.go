package cli

import (
	"path/filepath"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/export"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var flagPDF bool

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the finalized paper (--pdf requires pandoc + wkhtmltopdf)",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			finalDir := filepath.Join(state.ArtifactsDir(root), "final")

			if flagPDF {
				return export.ExportPDF(finalDir)
			}

			return cmd.Help()
		},
	}

	cmd.Flags().BoolVar(&flagPDF, "pdf", false, "export to PDF via pandoc + wkhtmltopdf")
	return cmd
}

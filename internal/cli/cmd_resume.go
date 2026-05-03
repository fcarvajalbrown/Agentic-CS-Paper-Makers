package cli

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

var stageNextCmd = map[state.Stage]string{
	state.StageInit:      "paperflow research",
	state.StageResearch:  "paperflow architect",
	state.StageArchitect: "paperflow write",
	state.StageWrite:     "paperflow review",
	state.StageReview:    "paperflow inbox",
	state.StageInbox:     "paperflow finalize",
	state.StageFinalize:  "paperflow export --pdf  (optional — paper is ready)",
	state.StageDone:      "(workflow complete)",
}

func newResumeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "resume",
		Short: "Print the next command to run after a crash or pause",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := requireProject()
			if err != nil {
				return err
			}

			cp, err := state.LoadCheckpoint(root)
			if err != nil {
				return err
			}

			stage := state.Stage(cp.Stage)
			if next, ok := stageNextCmd[stage]; ok {
				fmt.Printf("Current stage: %s\nNext command:  %s\n", cp.Stage, next)
			} else {
				fmt.Printf("Current stage: %s\nRun 'paperflow status' for details.\n", cp.Stage)
			}
			return nil
		},
	}
}

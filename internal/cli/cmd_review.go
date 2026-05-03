package cli

import (
	"encoding/json"
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

// criticResp is the shape each critic returns.
type criticResp struct {
	Reviewer string `json:"reviewer"`
	Content  string `json:"content"`
}

// metaReviewerResp is the shape the meta-reviewer returns.
type metaReviewerResp struct {
	Reviewer      string `json:"reviewer"`
	Contradictions []struct {
		Critics                []string `json:"critics"`
		Description            string   `json:"description"`
		HumanDecisionRequired  string   `json:"human_decision_required"`
	} `json:"contradictions"`
	Agreements []struct {
		Critics     []string `json:"critics"`
		Description string   `json:"description"`
	} `json:"agreements"`
	Summary string `json:"summary"`
}

func newReviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "Dispatch 3 critics in parallel, then run Meta-Reviewer",
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

			client, _, err := newClient(cfg, "formalist", state.CacheDir(root))
			if err != nil {
				return err
			}

			orch := agents.NewOrchestrator(agents.NewRunner(client))

			fmt.Println("Running 3 critics in parallel...")
			results, err := orch.RunCritics(string(draftData))
			if err != nil {
				return fmt.Errorf("critics: %w", err)
			}

			// Parse and validate each critic response.
			reviews := make(map[string]struct {
				Content string `json:"content"`
				Status  string `json:"status"`
			})

			for _, r := range results {
				var cr criticResp
				if parseErr := json.Unmarshal([]byte(r.Response), &cr); parseErr != nil {
					reviews[string(r.Agent)] = struct {
						Content string `json:"content"`
						Status  string `json:"status"`
					}{Content: r.Response, Status: "pending"}
				} else {
					reviews[string(r.Agent)] = struct {
						Content string `json:"content"`
						Status  string `json:"status"`
					}{Content: cr.Content, Status: "pending"}
				}
				fmt.Printf("  [%s] done\n", r.Agent)
			}

			// Run Meta-Reviewer with all 3 reviews.
			metaMsg, _ := json.Marshal(map[string]string{
				"formalist_review":     reviews["formalist"].Content,
				"game_theorist_review": reviews["game_theorist"].Content,
				"skeptic_review":       reviews["skeptic"].Content,
			})

			metaClient, _, err := newClient(cfg, "meta_reviewer", state.CacheDir(root))
			if err != nil {
				return err
			}
			metaRunner := agents.NewRunner(metaClient)
			metaResp, err := metaRunner.Run(agents.RunRequest{
				Agent:    agents.AgentMetaReviewer,
				Stage:    "meta_reviewer",
				Messages: []llm.Message{{Role: "user", Content: string(metaMsg)}},
			})
			if err != nil {
				return fmt.Errorf("meta-reviewer: %w", err)
			}

			// Build contradictions list as strings for the review_round schema.
			var contradictions []string
			var meta metaReviewerResp
			if err := json.Unmarshal([]byte(metaResp), &meta); err == nil {
				for _, c := range meta.Contradictions {
					contradictions = append(contradictions, c.Description)
				}
			}

			// Assemble review_round artifact.
			roundNum := 1
			if existing, err := artifacts.LatestVersion(artDir, "review_round"); err == nil {
				roundNum = existing + 1
			}

			type reviewEntry struct {
				Content string `json:"content"`
				Status  string `json:"status"`
			}
			type metaNotes struct {
				Contradictions []string `json:"contradictions"`
			}
			type reviewRound struct {
				Round          int                    `json:"round"`
				Reviews        map[string]reviewEntry `json:"reviews"`
				MetaReviewerNotes *metaNotes          `json:"meta_reviewer_notes,omitempty"`
			}

			round := reviewRound{
				Round: roundNum,
				Reviews: map[string]reviewEntry{
					"formalist":     {Content: reviews["formalist"].Content, Status: "pending"},
					"game_theorist": {Content: reviews["game_theorist"].Content, Status: "pending"},
					"skeptic":       {Content: reviews["skeptic"].Content, Status: "pending"},
				},
			}
			if len(contradictions) > 0 {
				round.MetaReviewerNotes = &metaNotes{Contradictions: contradictions}
			}

			roundData, err := json.MarshalIndent(round, "", "  ")
			if err != nil {
				return err
			}

			if err := artifacts.Write(artDir, "review_round", roundNum, artifacts.SchemaReviewRound, roundData); err != nil {
				return fmt.Errorf("save review_round: %w", err)
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageReview)
			cp.LatestArtifact = artifacts.ArtifactPath(artDir, "review_round", roundNum)
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			fmt.Printf("Review complete: review_round_v%d.json\n", roundNum)
			if len(contradictions) > 0 {
				fmt.Printf("Meta-Reviewer flagged %d conflict(s).\n", len(contradictions))
			}
			fmt.Println("Next step: paperflow inbox")
			return nil
		},
	}
}

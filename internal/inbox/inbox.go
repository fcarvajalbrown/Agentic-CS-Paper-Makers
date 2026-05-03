package inbox

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
)

type ReviewStatus string

const (
	Pending   ReviewStatus = "pending"
	Addressed ReviewStatus = "addressed"
	Skipped   ReviewStatus = "skipped"
)

type Review struct {
	Content string       `json:"content"`
	Status  ReviewStatus `json:"status"`
}

type HumanResponse struct {
	HumanSays             string `json:"human_says"`
	WriterRevisionApplied bool   `json:"writer_revision_applied"`
}

type MetaReviewerNotes struct {
	Contradictions []string `json:"contradictions"`
}

type ReviewRound struct {
	Round          int                      `json:"round"`
	Reviews        map[string]Review        `json:"reviews"`
	HumanResponses map[string]HumanResponse `json:"human_responses,omitempty"`
	MetaReviewer   *MetaReviewerNotes       `json:"meta_reviewer_notes,omitempty"`
}

type Inbox struct {
	runner       *agents.Runner
	artifactsDir string
	out          io.Writer
	in           *bufio.Scanner
}

func New(runner *agents.Runner, artifactsDir string) *Inbox {
	return &Inbox{
		runner:       runner,
		artifactsDir: artifactsDir,
		out:          os.Stdout,
		in:           bufio.NewScanner(os.Stdin),
	}
}

var reviewerOrder = []string{"formalist", "game_theorist", "skeptic"}

var reviewerLabel = map[string]string{
	"formalist":     "Formalist",
	"game_theorist": "Game Theorist",
	"skeptic":       "Skeptic",
}

// Run starts the interactive inbox loop for the given review round version.
// It saves the updated round after each action and returns when the user types "stop".
func (b *Inbox) Run(roundVersion int) error {
	var round ReviewRound
	if err := artifacts.ReadInto(b.artifactsDir, "review_round", roundVersion, &round); err != nil {
		return fmt.Errorf("inbox: load review_round: %w", err)
	}
	if round.HumanResponses == nil {
		round.HumanResponses = make(map[string]HumanResponse)
	}

	draftVersion, err := artifacts.LatestVersion(b.artifactsDir, "draft")
	if err != nil {
		return fmt.Errorf("inbox: load draft: %w", err)
	}

	for {
		b.printStatus(&round)

		fmt.Fprintf(b.out, "\nSelect reviewer (1-%d), 'skip <N>', or 'stop' to finalize: ", len(reviewerOrder))
		if !b.in.Scan() {
			break
		}
		input := strings.TrimSpace(b.in.Text())

		switch {
		case input == "stop" || input == "":
			fmt.Fprintln(b.out, "\nExiting inbox. Run 'paperflow finalize' when ready.")
			return nil

		case strings.HasPrefix(input, "skip "):
			idx := b.parseIndex(strings.TrimPrefix(input, "skip "))
			if idx < 0 {
				fmt.Fprintln(b.out, "Invalid reviewer number.")
				continue
			}
			name := reviewerOrder[idx]
			r := round.Reviews[name]
			r.Status = Skipped
			round.Reviews[name] = r
			if err := b.saveRound(roundVersion, &round); err != nil {
				return err
			}
			fmt.Fprintf(b.out, "[%s marked as SKIPPED]\n", reviewerLabel[name])

		default:
			idx := b.parseIndex(input)
			if idx < 0 {
				fmt.Fprintln(b.out, "Invalid selection. Enter a number, 'skip <N>', or 'stop'.")
				continue
			}
			name := reviewerOrder[idx]
			review := round.Reviews[name]

			fmt.Fprintf(b.out, "\n--- %s Review ---\n%s\n---\n", reviewerLabel[name], review.Content)
			fmt.Fprintf(b.out, "\nYour response to %s (blank line to finish):\n> ", reviewerLabel[name])

			var lines []string
			for b.in.Scan() {
				line := b.in.Text()
				if line == "" {
					break
				}
				lines = append(lines, line)
			}
			if len(lines) == 0 {
				fmt.Fprintln(b.out, "No response entered.")
				continue
			}
			humanSays := strings.Join(lines, "\n")

			fmt.Fprintln(b.out, "\n[Sending to Writer for micro-revision...]")

			draftData, err := artifacts.Read(b.artifactsDir, "draft", draftVersion)
			if err != nil {
				return err
			}

			userMsg := fmt.Sprintf(
				"Perform a targeted micro-revision of the following draft.\n\nCurrent draft (JSON):\n%s\n\n%s reviewer raised this critique:\n%s\n\nHuman researcher direction:\n%s\n\nRevise only the sections relevant to this critique. Output the complete revised draft as JSON in the same format.",
				string(draftData), reviewerLabel[name], review.Content, humanSays,
			)

			resp, err := b.runner.Run(agents.RunRequest{
				Agent:    agents.AgentWriter,
				Stage:    "writer",
				Messages: []llm.Message{{Role: "user", Content: userMsg}},
			})
			if err != nil {
				return fmt.Errorf("writer micro-revision: %w", err)
			}

			draftVersion++
			if err := artifacts.Write(b.artifactsDir, "draft", draftVersion, artifacts.SchemaDraft, []byte(resp)); err != nil {
				return fmt.Errorf("save revised draft v%d: %w", draftVersion, err)
			}

			review.Status = Addressed
			round.Reviews[name] = review
			round.HumanResponses[name] = HumanResponse{
				HumanSays:             humanSays,
				WriterRevisionApplied: true,
			}
			if err := b.saveRound(roundVersion, &round); err != nil {
				return err
			}

			fmt.Fprintf(b.out, "[Draft v%d generated. %s marked as ADDRESSED]\n", draftVersion, reviewerLabel[name])
		}
	}

	return nil
}

func (b *Inbox) printStatus(round *ReviewRound) {
	fmt.Fprintf(b.out, "\nReview Round %d — %d reviews\n", round.Round, len(reviewerOrder))
	for i, name := range reviewerOrder {
		r := round.Reviews[name]
		fmt.Fprintf(b.out, "[%d] %-16s — %s\n", i+1, reviewerLabel[name], strings.ToUpper(string(r.Status)))
	}
	if round.MetaReviewer != nil && len(round.MetaReviewer.Contradictions) > 0 {
		fmt.Fprintln(b.out, "\nConflicts flagged by Meta-Reviewer:")
		for _, c := range round.MetaReviewer.Contradictions {
			fmt.Fprintf(b.out, "  [!] %s\n", c)
		}
	}
}

func (b *Inbox) parseIndex(s string) int {
	s = strings.TrimSpace(s)
	var n int
	if _, err := fmt.Sscan(s, &n); err != nil {
		return -1
	}
	if n < 1 || n > len(reviewerOrder) {
		return -1
	}
	return n - 1
}

func (b *Inbox) saveRound(version int, round *ReviewRound) error {
	data, err := json.MarshalIndent(round, "", "  ")
	if err != nil {
		return err
	}
	return artifacts.Write(b.artifactsDir, "review_round", version, artifacts.SchemaReviewRound, data)
}

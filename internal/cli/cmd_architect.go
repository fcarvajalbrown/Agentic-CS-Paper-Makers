package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/agents"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/llm"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

// architectResp is the envelope the Architect agent uses for all responses.
type architectResp struct {
	Type string          `json:"type"`
	Text string          `json:"text,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

func newArchitectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "architect",
		Short: "Run the Socratic Architect (human gate: approve blueprint before proceeding)",
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
			rcVer, err := artifacts.LatestVersion(artDir, "research_context")
			if err != nil {
				return fmt.Errorf("no research_context found; run 'paperflow research' first")
			}
			rcData, err := artifacts.Read(artDir, "research_context", rcVer)
			if err != nil {
				return err
			}

			client, _, err := newClient(cfg, "architect", state.CacheDir(root))
			if err != nil {
				return err
			}

			runner := agents.NewRunner(client)
			scanner := bufio.NewScanner(os.Stdin)

			seed := cfg.Seed
			userMsg := fmt.Sprintf("Research context:\n%s\n\nSeed: %s", string(rcData), seed)
			messages := []llm.Message{{Role: "user", Content: userMsg}}

			var blueprintData json.RawMessage

			for {
				resp, err := runner.Run(agents.RunRequest{
					Agent:    agents.AgentArchitect,
					Stage:    "architect",
					Messages: messages,
				})
				if err != nil {
					return fmt.Errorf("architect agent: %w", err)
				}

				// Append assistant response to history.
				messages = append(messages, llm.Message{Role: "assistant", Content: resp})

				var env architectResp
				if err := json.Unmarshal([]byte(resp), &env); err != nil {
					return fmt.Errorf("architect returned unexpected format: %w\nRaw: %s", err, resp)
				}

				switch env.Type {
				case "question":
					fmt.Printf("\nArchitect: %s\n> ", env.Text)
					if !scanner.Scan() {
						return fmt.Errorf("input closed before blueprint was produced")
					}
					answer := strings.TrimSpace(scanner.Text())
					messages = append(messages, llm.Message{Role: "user", Content: answer})

				case "blueprint":
					blueprintData = env.Data
					goto approve

				default:
					return fmt.Errorf("unknown architect response type: %q", env.Type)
				}
			}

		approve:
			fmt.Println("\n--- Proposed Blueprint ---")
			pretty, _ := json.MarshalIndent(blueprintData, "", "  ")
			fmt.Println(string(pretty))
			fmt.Print("\nApprove this blueprint? [y/n]: ")

			if !scanner.Scan() {
				return fmt.Errorf("input closed")
			}
			if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
				fmt.Println("Blueprint rejected. Re-run 'paperflow architect' to start over.")
				return nil
			}

			// Inject human_approved = true and save.
			var bpMap map[string]any
			if err := json.Unmarshal(blueprintData, &bpMap); err != nil {
				return err
			}
			bpMap["human_approved"] = true
			bpMap["research_context_id"] = fmt.Sprintf("research_context_v%d", rcVer)

			finalBlueprint, err := json.MarshalIndent(bpMap, "", "  ")
			if err != nil {
				return err
			}

			ver := 1
			if err := artifacts.Write(artDir, "blueprint", ver, artifacts.SchemaBlueprint, finalBlueprint); err != nil {
				return fmt.Errorf("save blueprint: %w", err)
			}

			cp, _ := state.LoadCheckpoint(root)
			cp.Stage = string(state.StageArchitect)
			cp.LatestArtifact = artifacts.ArtifactPath(artDir, "blueprint", ver)
			if err := state.SaveCheckpoint(root, cp); err != nil {
				return err
			}

			fmt.Printf("\nBlueprint approved and saved: blueprint_v%d.json\n", ver)
			fmt.Println("Next step: paperflow write")
			return nil
		},
	}
}

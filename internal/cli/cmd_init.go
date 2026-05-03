package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/state"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init <paper-name>",
		Short: "Create a new paper project directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if err := os.Mkdir(name, 0755); err != nil {
				return fmt.Errorf("init: %w", err)
			}

			if err := state.InitProject(name); err != nil {
				return fmt.Errorf("init: scaffold: %w", err)
			}

			seed, err := promptSeed()
			if err != nil {
				return err
			}

			projectCfg := config.FileConfig{Seed: seed}
			cfgData, err := json.MarshalIndent(projectCfg, "", "  ")
			if err != nil {
				return err
			}
			cfgPath := filepath.Join(name, config.ConfigFile)
			if err := os.WriteFile(cfgPath, cfgData, 0644); err != nil {
				return fmt.Errorf("init: write config: %w", err)
			}

			cp := &state.Checkpoint{Stage: string(state.StageInit)}
			if err := state.SaveCheckpoint(name, cp); err != nil {
				return fmt.Errorf("init: save checkpoint: %w", err)
			}

			fmt.Printf("Project created: %s/\n", name)
			fmt.Printf("API key: set %s in your environment.\n", config.EnvAPIKey)
			fmt.Printf("Next step: cd %s && paperflow research\n", name)
			return nil
		},
	}
}

func promptSeed() (string, error) {
	fmt.Print("Seed idea for your paper (describe the research question): ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		seed := strings.TrimSpace(scanner.Text())
		if seed != "" {
			return seed, nil
		}
	}
	return "", fmt.Errorf("seed cannot be empty")
}

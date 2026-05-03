package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/artifacts"
	"github.com/spf13/cobra"
)

// schemaByFilename maps filename substrings to schema names.
var schemaByFilename = []struct {
	contains string
	schema   artifacts.SchemaName
}{
	{"research_context", artifacts.SchemaResearchContext},
	{"blueprint", artifacts.SchemaBlueprint},
	{"draft", artifacts.SchemaDraft},
	{"review_round", artifacts.SchemaReviewRound},
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file.json>",
		Short: "Validate an artifact JSON file against its schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("validate: read file: %w", err)
			}

			schema := detectSchema(path)
			if schema == "" {
				return fmt.Errorf("cannot detect schema for %q — filename must contain one of: research_context, blueprint, draft, review_round", path)
			}

			if err := artifacts.Validate(schema, data); err != nil {
				fmt.Fprintf(os.Stderr, "INVALID: %v\n", err)
				return fmt.Errorf("validation failed")
			}

			fmt.Printf("VALID: %s passes %s\n", path, schema)
			return nil
		},
	}
}

func detectSchema(path string) artifacts.SchemaName {
	lower := strings.ToLower(path)
	for _, entry := range schemaByFilename {
		if strings.Contains(lower, entry.contains) {
			return entry.schema
		}
	}
	return ""
}

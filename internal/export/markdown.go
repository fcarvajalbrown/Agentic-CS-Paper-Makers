package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FinalizeOutput matches the JSON the Finalize Agent returns.
type FinalizeOutput struct {
	PaperMarkdown string        `json:"paper_markdown"`
	ReferencesBib string        `json:"references_bib"`
	Metadata      FinalMetadata `json:"metadata"`
}

// FinalMetadata holds paper-level statistics written to metadata.json.
type FinalMetadata struct {
	Title         string   `json:"title"`
	Date          string   `json:"date"`
	Keywords      []string `json:"keywords"`
	CitationCount int      `json:"citation_count"`
	SectionCount  int      `json:"section_count"`
	HasTheorems   bool     `json:"has_theorems"`
	HasLemmas     bool     `json:"has_lemmas"`
}

// WriteArtifacts writes paper.md, references.bib, and metadata.json to finalDir.
func WriteArtifacts(finalDir string, output *FinalizeOutput) error {
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		return fmt.Errorf("export: create final dir: %w", err)
	}

	if err := writeFile(filepath.Join(finalDir, "paper.md"), []byte(output.PaperMarkdown)); err != nil {
		return fmt.Errorf("export: write paper.md: %w", err)
	}

	if err := writeFile(filepath.Join(finalDir, "references.bib"), []byte(output.ReferencesBib)); err != nil {
		return fmt.Errorf("export: write references.bib: %w", err)
	}

	metaJSON, err := json.MarshalIndent(output.Metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("export: marshal metadata: %w", err)
	}
	if err := writeFile(filepath.Join(finalDir, "metadata.json"), metaJSON); err != nil {
		return fmt.Errorf("export: write metadata.json: %w", err)
	}

	return nil
}

func writeFile(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

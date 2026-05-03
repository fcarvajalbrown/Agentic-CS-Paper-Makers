package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Checkpoint struct {
	Stage         string  `json:"current_stage"`
	LatestArtifact string  `json:"latest_artifact"`
	TotalTokens   int     `json:"total_tokens"`
	TotalCost     float64 `json:"total_cost"`
}

func LoadCheckpoint(root string) (*Checkpoint, error) {
	path := filepath.Join(root, ".paperflow", "state.json")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Checkpoint{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var cp Checkpoint
	if err := json.NewDecoder(f).Decode(&cp); err != nil {
		return nil, err
	}
	return &cp, nil
}

func SaveCheckpoint(root string, cp *Checkpoint) error {
	path := filepath.Join(root, ".paperflow", "state.json")

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cp); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, path)
}

func RollbackTo(root, stage string) error {
	cp, err := LoadCheckpoint(root)
	if err != nil {
		return err
	}
	cp.Stage = stage
	cp.LatestArtifact = ""
	return SaveCheckpoint(root, cp)
}

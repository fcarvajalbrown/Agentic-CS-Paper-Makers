package state

import (
	"os"
	"path/filepath"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/config"
)

var projectDirs = []string{
	config.ArtifactsDir,
	config.ArtifactsDir + "/final",
	config.CacheDir,
	config.ErrorsDir,
}

func InitProject(root string) error {
	for _, dir := range projectDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			return err
		}
	}
	return nil
}

func ProjectRoot() (string, error) {
	return os.Getwd()
}

func PaperflowDir(root string) string {
	return filepath.Join(root, ".paperflow")
}

func ArtifactsDir(root string) string {
	return filepath.Join(root, config.ArtifactsDir)
}

func CacheDir(root string) string {
	return filepath.Join(root, config.CacheDir)
}

func ErrorsDir(root string) string {
	return filepath.Join(root, config.ErrorsDir)
}

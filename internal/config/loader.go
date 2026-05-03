package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type FileConfig struct {
	APIKey    string            `json:"api_key,omitempty"`
	Model     string            `json:"model,omitempty"`
	Budget    float64           `json:"budget,omitempty"`
	MaxTokens int               `json:"max_tokens,omitempty"`
	Agents    map[string]string `json:"agents,omitempty"`
}

func LoadProjectConfig(projectDir string) (*FileConfig, error) {
	return loadJSON(filepath.Join(projectDir, ConfigFile))
}

func LoadGlobalConfig() (*FileConfig, error) {
	return loadJSON(globalConfigPath())
}

func SaveProjectConfig(projectDir string, cfg *FileConfig) error {
	return saveJSON(filepath.Join(projectDir, ConfigFile), cfg)
}

func globalConfigPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), GlobalConfigDir, GlobalConfigFile)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", GlobalConfigDir, GlobalConfigFile)
}

func loadJSON(path string) (*FileConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileConfig{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg FileConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveJSON(path string, cfg *FileConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
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

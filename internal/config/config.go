package config

import "os"

type Config struct {
	APIKey    string
	Model     string
	Budget    float64
	MaxTokens int
	Agents    map[string]string
	Seed      string
}

func Load(projectDir string, flags Flags) (*Config, error) {
	global, err := LoadGlobalConfig()
	if err != nil {
		return nil, err
	}

	project, err := LoadProjectConfig(projectDir)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		APIKey:    "",
		Model:     DefaultModel,
		Budget:    DefaultBudget,
		MaxTokens: DefaultMaxTokens,
		Agents:    make(map[string]string),
	}

	applyFile(cfg, global)
	applyFile(cfg, project)
	applyEnv(cfg)
	applyFlags(cfg, flags)

	return cfg, nil
}

func (c *Config) ModelForAgent(agent string) string {
	if m, ok := c.Agents[agent]; ok {
		return m
	}
	if m, ok := ModelByStage[agent]; ok {
		return m
	}
	return c.Model
}

type Flags struct {
	Model     string
	Budget    float64
	MaxTokens int
	Cheap     bool
	Production bool
	AgentModel map[string]string
}

func applyFile(cfg *Config, f *FileConfig) {
	if f.APIKey != "" {
		cfg.APIKey = f.APIKey
	}
	if f.Model != "" {
		cfg.Model = f.Model
	}
	if f.Budget != 0 {
		cfg.Budget = f.Budget
	}
	if f.MaxTokens != 0 {
		cfg.MaxTokens = f.MaxTokens
	}
	for k, v := range f.Agents {
		cfg.Agents[k] = v
	}
	if f.Seed != "" {
		cfg.Seed = f.Seed
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv(EnvAPIKey); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv(EnvModel); v != "" {
		cfg.Model = v
	}
}

func applyFlags(cfg *Config, f Flags) {
	if f.Cheap {
		cfg.Model = CheapModel
		for k := range ModelByStage {
			cfg.Agents[k] = CheapModel
		}
	}
	if f.Production {
		cfg.Model = ProductionModel
		for k := range ModelByStage {
			cfg.Agents[k] = ProductionModel
		}
	}
	if f.Model != "" {
		cfg.Model = f.Model
	}
	if f.Budget != 0 {
		cfg.Budget = f.Budget
	}
	if f.MaxTokens != 0 {
		cfg.MaxTokens = f.MaxTokens
	}
	for k, v := range f.AgentModel {
		cfg.Agents[k] = v
	}
}

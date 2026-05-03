package agents

import (
	"fmt"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/pkg/embed"
)

type AgentName string

const (
	AgentResearch      AgentName = "research"
	AgentArchitect     AgentName = "architect"
	AgentWriter        AgentName = "writer"
	AgentFormalist     AgentName = "formalist"
	AgentGameTheorist  AgentName = "game_theorist"
	AgentSkeptic       AgentName = "skeptic"
	AgentMetaReviewer  AgentName = "meta_reviewer"
	AgentFinalize      AgentName = "finalize"
)

// profileFile maps each agent name to its embedded .md filename.
var profileFile = map[AgentName]string{
	AgentResearch:     "agents/research.md",
	AgentArchitect:    "agents/architect.md",
	AgentWriter:       "agents/writer.md",
	AgentFormalist:    "agents/formalist.md",
	AgentGameTheorist: "agents/game_theorist.md",
	AgentSkeptic:      "agents/skeptic.md",
	AgentMetaReviewer: "agents/meta_reviewer.md",
	AgentFinalize:     "agents/finalize.md",
}

// LoadProfile returns the system prompt for the named agent.
func LoadProfile(name AgentName) (string, error) {
	path, ok := profileFile[name]
	if !ok {
		return "", fmt.Errorf("unknown agent: %q", name)
	}
	data, err := embed.AgentProfiles.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load profile %q: %w", name, err)
	}
	return string(data), nil
}

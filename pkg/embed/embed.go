package embed

import "embed"

//go:embed agents/*.md
var AgentProfiles embed.FS

//go:embed schemas/*.json
var Schemas embed.FS

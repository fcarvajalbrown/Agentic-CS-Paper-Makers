package artifacts

import "github.com/fcarvajalbrown/agentic-cs-paper-makers/pkg/embed"

type SchemaName string

const (
	SchemaResearchContext SchemaName = "research_context.schema.json"
	SchemaBlueprint       SchemaName = "blueprint.schema.json"
	SchemaDraft           SchemaName = "draft.schema.json"
	SchemaReviewRound     SchemaName = "review_round.schema.json"
)

func LoadSchema(name SchemaName) ([]byte, error) {
	return embed.Schemas.ReadFile("schemas/" + string(name))
}

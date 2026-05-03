package config

const (
	DefaultModel    = "kimi-k2.5"
	CheapModel      = "moonshot-v1-32k"
	ProductionModel = "kimi-k2.6"

	DefaultBudget    = 10.00
	DefaultMaxTokens = 32000

	ArtifactsDir = ".paperflow/artifacts"
	CacheDir     = ".paperflow/cache"
	ErrorsDir    = ".paperflow/errors"
	StateFile    = ".paperflow/state.json"
	ConfigFile   = ".paperflow/config.json"
	LogFile      = "paperflow.log"

	GlobalConfigDir  = "paperflow"
	GlobalConfigFile = "config.json"

	EnvAPIKey = "PAPERFLOW_API_KEY"
	EnvModel  = "PAPERFLOW_MODEL"
)

var ModelByStage = map[string]string{
	"research":      CheapModel,
	"architect":     DefaultModel,
	"writer":        DefaultModel,
	"formalist":     DefaultModel,
	"game_theorist": DefaultModel,
	"skeptic":       DefaultModel,
	"meta_reviewer": DefaultModel,
	"finalize":      CheapModel,
}

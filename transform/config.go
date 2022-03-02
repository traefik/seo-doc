package transform

import "github.com/urfave/cli/v2"

// Transform flag names.
const (
	FlagPath    = "path"
	FlagProduct = "product"
)

// Config is the bot configuration.
type Config struct {
	Path    string
	Product string
}

// NewConfig creates a new Config.
func NewConfig(cliCtx *cli.Context) Config {
	return Config{
		Path:    cliCtx.Path(FlagPath),
		Product: cliCtx.String(FlagProduct),
	}
}

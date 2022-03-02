package sitemap

import (
	"github.com/urfave/cli/v2"
)

const flagRoot = "root"

// Command is the sitemap command.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "sitemap",
		Usage:       "Generates sitemap of the documentation.",
		Description: "Sitemap generator.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  flagRoot,
				Usage: "Path to the root of the documentation.",
				Value: ".",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			return Generate(cliCtx.Path(flagRoot))
		},
	}
}

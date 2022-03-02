package sitemap

import (
	"github.com/ettle/strcase"
	"github.com/urfave/cli/v2"
)

const (
	flagRoot         = "root"
	flagDebug        = "debug"
	flagGitUserName  = "git-user-name"
	flagGitUserEmail = "git-user-email"
	flagGithubToken  = "token"
)

// Command is the sitemap command.
func Command() *cli.Command {
	return &cli.Command{
		Name:        "sitemap",
		Usage:       "Generates sitemap of the documentation.",
		Description: "Sitemap generator.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   flagRoot,
				Usage:  "Path to the root of the documentation.",
				Value:  ".",
				Hidden: true,
			},
			&cli.BoolFlag{
				Name:    flagDebug,
				Usage:   "Path to the root of the documentation.",
				EnvVars: []string{strcase.ToSNAKE(flagDebug)},
			},
			&cli.StringFlag{
				Name:     flagGitUserName,
				Usage:    "UserName used to commit the sitemap files.",
				EnvVars:  []string{strcase.ToSNAKE(flagGitUserName)},
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagGitUserEmail,
				Usage:    "Email used to commit the sitemap files.",
				EnvVars:  []string{strcase.ToSNAKE(flagGitUserEmail)},
				Required: true,
			},
			&cli.StringFlag{
				Name:     flagGithubToken,
				Usage:    "GitHub token.",
				EnvVars:  []string{"GITHUB_TOKEN"},
				Required: true,
			},
		},
		Action: func(cliCtx *cli.Context) error {
			err := Generate(cliCtx.Path(flagRoot))
			if err != nil {
				return err
			}

			return Commit(NewGitInfo(cliCtx), cliCtx.Bool(flagDebug))
		},
	}
}

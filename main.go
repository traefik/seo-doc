package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/traefik/seo/transform"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "seo",
		Description: "Documentation modification for SEO.",
		Usage:       "SEO doc",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  transform.FlagPath,
				Usage: "Path of the documentation.",
			},
			&cli.StringFlag{
				Name:  transform.FlagProduct,
				Usage: "Product name.",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			config := transform.NewConfig(cliCtx)

			err := validate(config)
			if err != nil {
				return err
			}

			return transform.Run(config)
		},
		Commands: []*cli.Command{
			{
				Name:        "version",
				Usage:       "Version information.",
				Description: "Version information.",
				Action: func(cliCtx *cli.Context) error {
					displayVersion()
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Error while executing command ", err)
	}
}

func validate(cfg transform.Config) error {
	if strings.TrimSpace(cfg.Path) == "" {
		return errors.New("path is required")
	}

	return nil
}

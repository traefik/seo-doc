package transform

import (
	"os"
	"path/filepath"
)

type fileTransform interface {
	Match(path string) bool
	Apply(path string) error
}

// Run applies transformations is needed.
func Run(cfg Config) error {
	productName := getProductName(cfg)

	transforms := []fileTransform{
		NewPageTransform(productName),
		NewSitemapTransform(productName),
	}

	return filepath.Walk(cfg.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, transform := range transforms {
				if transform.Match(path) {
					return transform.Apply(path)
				}
			}

			return nil
		},
	)
}

func getProductName(cfg Config) string {
	if cfg.Product != "" {
		return cfg.Product
	}

	return filepath.Base(cfg.Path)
}

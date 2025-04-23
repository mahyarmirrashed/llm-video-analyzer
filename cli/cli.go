package cli

import (
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func New() *cli.App {
	cfg := &config.Config{}

	app := &cli.App{
		Name:  "llm-video-analyze",
		Usage: "Search through videos using natural language",
		Commands: []*cli.Command{
			ParseCommand(cfg),
			QueryCommand(cfg),
			CleanCommand(cfg),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "ollama-url",
				Value:       "http://localhost:11434",
				Usage:       "Ollama server URL",
				Destination: &cfg.OllamaURL,
			},
			&cli.StringFlag{
				Name:        "database-url",
				Value:       "http://localhost:6334",
				Usage:       "Vector database URL",
				Destination: &cfg.DatabaseURL,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "Enable debug mode",
				Destination: &cfg.Debug,
			},
		},
	}

	return app
}

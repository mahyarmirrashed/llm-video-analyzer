package cli

import (
	"github.com/urfave/cli/v2"
)

type Config struct {
	SamplingInterval int
	SamplingModel    string
	OllamaURL        string
	DatabaseURL      string
	Debug            bool
}

func New() *cli.App {
	cfg := &Config{}

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
				Value:       "http://localhost:6333",
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

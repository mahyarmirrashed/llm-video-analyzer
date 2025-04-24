package cli

import (
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func QueryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "Query processed videos",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "query-model",
				Value:       "llama3.2",
				Usage:       "Query model for search",
				Destination: &cfg.QueryModel,
			},
			&cli.IntFlag{
				Name:        "query-limit",
				Value:       3,
				Usage:       "Number of results to return",
				Destination: &cfg.QueryLimit,
			},
		},
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

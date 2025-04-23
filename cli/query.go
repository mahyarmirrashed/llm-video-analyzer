package cli

import (
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func QueryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "Search processed videos",
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

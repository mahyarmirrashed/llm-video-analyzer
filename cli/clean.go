package cli

import (
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func CleanCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "Clean processed video from database",
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

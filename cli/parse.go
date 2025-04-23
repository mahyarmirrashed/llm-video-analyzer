package cli

import (
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func ParseCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "parse",
		Usage: "Process a video into database",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "sampling-interval",
				Value:       2,
				Usage:       "Frame sampling interval (seconds)",
				Destination: &cfg.SamplingInterval,
			},
			&cli.StringFlag{
				Name:        "sampling-model",
				Value:       "llava",
				Usage:       "Frame sampling model for analysis",
				Destination: &cfg.SamplingModel,
			},
		},
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

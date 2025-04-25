package cli

import (
	"fmt"
	"log"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/cmd"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/urfave/cli/v2"
)

func ProcessCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "process",
		ArgsUsage: "<youtube-url>",
		Usage:     "Process a video into database",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "sampling-interval",
				Value:       2,
				Usage:       "Frame sampling interval (seconds)",
				Destination: &cfg.SamplingInterval,
			},
			&cli.StringFlag{
				Name:        "sampling-model",
				Value:       "llava:7b",
				Usage:       "Frame sampling model for analysis",
				Destination: &cfg.SamplingModel,
			},
		},
		Action: func(c *cli.Context) error {
			url := c.Args().First()
			if url == "" {
				return fmt.Errorf("youtube url is required")
			}

			db, err := qdrant.New(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database")
			}

			command := cmd.New(cfg, db)
			_, err = command.Process(c.Context, url)
			if err != nil {
				return err
			}

			log.Println("successfully processed video")

			return nil
		},
	}
}

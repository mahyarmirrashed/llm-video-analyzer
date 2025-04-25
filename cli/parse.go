package cli

import (
	"fmt"
	"log"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/video"
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
			videoPath := c.Args().First()
			if videoPath == "" {
				return fmt.Errorf("video path required")
			}

			v, err := video.New(videoPath)
			if err != nil {
				fmt.Errorf("failed to initialize video: %w", err)
			}

			if err := v.Extract(cfg.SamplingInterval); err != nil {
				return fmt.Errorf("frame extraction failed: %w", err)
			}

			defer func() {
				if !cfg.Debug {
					v.Cleanup()
				} else {
					log.Printf("temporary files retained at: %s", v.ProcessingPath)
				}
			}()

			dbClient, err := qdrant.New(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			for i := range v.Frames {
				log.Printf("processing frame %d...", i)

				frame := &v.Frames[i]

				if err := frame.Process(c.Context, cfg); err != nil {
					log.Printf("skipping frame %s: %v", frame.Path, err)
					continue
				}

				if err := dbClient.Store(c.Context, v.ID, frame); err != nil {
					log.Printf("failed to store frame %s: %v", frame.Path, err)
				}
			}

			log.Println("finished processing video")

			return nil
		},
	}
}

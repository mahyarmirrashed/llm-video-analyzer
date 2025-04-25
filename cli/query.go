package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/cmd"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/urfave/cli/v2"
)

func QueryCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "query",
		Usage:     "Query processed videos",
		ArgsUsage: "<query>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "query-model",
				Value:       "llama3.2",
				Usage:       "Query model for search",
				Destination: &cfg.QueryModel,
			},
			&cli.IntFlag{
				Name:        "limit",
				Value:       3,
				Usage:       "Number of results to return",
				Destination: &cfg.QueryLimit,
			},
		},
		Action: func(c *cli.Context) error {
			query := c.Args().First()
			if query == "" {
				return fmt.Errorf("search query required")
			}

			db, err := qdrant.New(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database")
			}

			command := cmd.New(cfg, db)
			pts, err := command.Query(c.Context, query, cfg.QueryLimit)
			if err != nil {
				return err
			}

			for i, res := range pts {
				parts := strings.SplitN(res.VideoID, "-", 2)

				filename := parts[1]

				log.Printf("Result: %d\n", i+1)
				log.Printf("  Video: %s\n", filename)
				log.Printf("  Command: mpv '%s' --start=%.0f\n", filename, res.Timestamp)
			}

			return nil
		},
	}
}

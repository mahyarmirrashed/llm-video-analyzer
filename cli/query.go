package cli

import (
	"fmt"
	"log"

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
				log.Printf("Result: %d\n", i+1)
				log.Printf("  Video: %s&t=%.0f\n", res.Url, res.Timestamp)
				log.Println()
			}

			return nil
		},
	}
}

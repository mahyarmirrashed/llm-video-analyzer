package cli

import (
	"fmt"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/ollama"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
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

			desc, err := ollama.GetDescriptionFromQuery(c.Context, cfg, query)
			if err != nil {
				return fmt.Errorf("failed to get description: %w", err)
			}

			embedding, err := ollama.GetTextEmbedding(c.Context, cfg, desc)
			if err != nil {
				return fmt.Errorf("failed to get embedding: %w", err)
			}

			dbClient, err := qdrant.New(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			pts, err := dbClient.Search(c.Context, embedding, uint64(cfg.QueryLimit))
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			} else if len(pts) == 0 {
				return fmt.Errorf("no results found")
			}

			for _, res := range pts {
				fmt.Printf("Result: %v", res)
			}

			return nil
		},
	}
}

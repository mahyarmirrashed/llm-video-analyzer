package cli

import (
	"fmt"

	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/qdrant"
	"github.com/urfave/cli/v2"
)

func CleanCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "Clean processed video from database",
		Action: func(c *cli.Context) error {
			dbClient, err := qdrant.New(cfg.DatabaseURL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			return dbClient.Cleanup(c.Context)
		},
	}
}

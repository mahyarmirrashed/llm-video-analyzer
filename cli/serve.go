package cli

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mahyarmirrashed/llm-video-analyzer/api"
	"github.com/mahyarmirrashed/llm-video-analyzer/pkg/config"
	"github.com/urfave/cli/v2"
)

func ServeCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Start the API server",
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:        "port",
				Value:       8080,
				Usage:       "Port to listen on",
				Destination: &cfg.ServerPort,
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}

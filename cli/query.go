package cli

import "github.com/urfave/cli/v2"

func QueryCommand(cfg *Config) *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "Search processed videos",
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

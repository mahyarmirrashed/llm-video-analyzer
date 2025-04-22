package cli

import "github.com/urfave/cli/v2"

func CleanCommand(cfg *Config) *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "Clean processed video from database",
		Action: func(c *cli.Context) error {
			return nil // TODO
		},
	}
}

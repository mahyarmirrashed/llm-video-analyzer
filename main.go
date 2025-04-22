package main

import (
	"log"
	"os"

	"github.com/mahyarmirrashed/llm-video-analyzer/cli"
)

func main() {
	app := cli.New()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

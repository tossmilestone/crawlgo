package main

import (
	"os"
	"fmt"

	"github.com/tossmilestone/crawlgo/cmd/crawlgo/app"
)

func main() {
	cmd := app.NewCrawlServerCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

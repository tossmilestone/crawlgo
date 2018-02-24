package app

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/cobra"

	"github.com/tossmilestone/crawlgo/pkg/server"
)

type Options struct {
	config *server.Config
}

func NewOptions() *Options {
	return &Options {
		config: new(server.Config),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.config.Site, "site", "", "The site to crawl.")
	fs.StringVar(&o.config.SaveDir, "save-dir", "./crawlgo", "The directory to save downloaded files.")
	fs.IntVar(&o.config.Workers, "workers", 3, "The number of workers to run the crawl tasks.")
	fs.StringVar(&o.config.DownloadSelector, "download-selector", "", "The DOM selector to query the links that will be downloaded from the site.")
}

func (o *Options) Run() {
	crawler, err := server.NewCrawler(o.config)
	if err != nil {
		log.Fatal(err)
		return
	}
	crawler.Run(make(chan struct{}))
}

func NewCrawlServerCommand() *cobra.Command {
	opts := NewOptions()
	cmd := &cobra.Command{
		Use: "crawlgo",
		Long: `The crawlgo is a server to crawl web sites in high concurrency written in golang.`,
		Run: func(cmd *cobra.Command, args[]string) {
			opts.Run()
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

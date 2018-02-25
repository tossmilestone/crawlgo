package app

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tossmilestone/crawlgo/pkg/server"
	"github.com/tossmilestone/crawlgo/pkg/version"
)

// Options represents the options to run the crawler server.
type Options struct {
	config *server.Config
}

// NewOptions creates a new Options object.
func NewOptions() *Options {
	return &Options{
		config: new(server.Config),
	}
}

// AddFlags adds flags to the crawlgo command.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.config.Site, "site", "", "The site to crawl.")
	fs.StringVar(&o.config.SaveDir, "save-dir", "./crawlgo", "The directory to save downloaded files.")
	fs.IntVar(&o.config.Workers, "workers", 0, "The number of workers to run the crawl tasks. If no set, will be 'runtime.NumCPU()'.")
	fs.StringVar(&o.config.DownloadSelector, "download-selector", "", "The DOM selector to query the links that will be downloaded from the site.")
}

// Run runs the crawler server.
func (o *Options) Run() {
	crawler, err := server.NewCrawler(o.config)
	if err != nil {
		log.Fatal(err)
		return
	}
	crawler.Run(make(chan struct{}))
}

// NewCrawlServerCommand creates a cobra.Command object to run the crawler server.
func NewCrawlServerCommand() *cobra.Command {
	opts := NewOptions()
	cmd := &cobra.Command{
		Version: version.VERSION,
		Use:  "crawlgo",
		Long: `The crawlgo is a server to crawl web sites in high concurrency written in golang.`,
		Run: func(cmd *cobra.Command, args []string) {
			opts.Run()
		},
	}
	opts.AddFlags(cmd.Flags())
	return cmd
}

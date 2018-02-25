package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/tossmilestone/crawlgo/pkg/js"
	"github.com/tossmilestone/crawlgo/pkg/util"
	"time"
)

// Config stores the configuration of the Crawler.
type Config struct {
	SaveDir          string
	Site             string
	Workers          int
	DownloadSelector string
}

// Crawler describes a Crawler server.
type Crawler struct {
	config     *Config
	phantom    *js.Phantom
	total      int
	failed     int
	downloaded int
	progress   chan int
}

// NewCrawler creates a Crawler object.
func NewCrawler(config *Config) (*Crawler, error) {
	// Create the save dir
	if config.SaveDir == "" {
		config.SaveDir = "./crawlgo"
	}
	err := os.MkdirAll(config.SaveDir, os.ModeDir)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Init workers
	if config.Workers == 0 {
		config.Workers = runtime.NumCPU()
	}

	return &Crawler{
		config:     config,
		phantom:    js.NewPhantom(),
		total:      0,
		downloaded: 0,
		progress:   make(chan int),
	}, nil
}

// Run runs a crawler to execute the crawl tasks.
func (c *Crawler) Run(stop chan struct{}) {
	log.Print("Running crawler...")

	go func() {
		c.startCrawl()
		log.Printf("All downloaded, exit.")
		stop <- struct{}{}
	}()

	go func() {
		c.displayProgress(stop)
	}()

	defer c.phantom.Stop()
	<-stop

	log.Print("Crawler stopped")
}

func (c *Crawler) startCrawl() error {
	// Site analysis
	c.phantom.Run()
	log.Printf("Opening %s ...", c.config.Site)
	links, err := c.phantom.ExtractLinksFromSelector(c.config.Site, c.config.DownloadSelector)
	if err != nil {
		return err
	}
	c.total = len(links)
	log.Printf("Downloading %d links", c.total)
	c.parallelizeDownload(links)
	return nil
}

func (c *Crawler) displayProgress(stop chan struct{}) {
	for {
		select {
		case result := <-c.progress:
			if result != 0 {
				c.failed++
			} else {
				c.downloaded++
			}
			log.Printf("Downloaded %d, failed %d links of total %d", c.downloaded, c.failed, c.total)
		case <-stop:
			return
		}
	}
}

func (c *Crawler) parallelizeDownload(downloads []interface{}) {
	downloadFile := func(index int) {
		log.Printf("Start to download %s", downloads[index])
		err := c.download(downloads[index].(string))
		if err != nil {
			log.Print(err)
			c.progress <- 1
		} else {
			c.progress <- 0
		}
	}
	util.Parallelize(c.config.Workers, len(downloads), downloadFile)
}

func (c *Crawler) download(urlStr string) error {
	frags := strings.Split(urlStr, "/")
	fileName := c.config.SaveDir + "/" + frags[len(frags)-1]

	file, err := os.Stat(fileName)
	if file != nil {
		return fmt.Errorf("%s exist", fileName)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(`download failed of "%s"`, frags[len(frags)-1])
	}

	out, err := os.Create(fileName)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Printf(`Downloaded of "%s"`, frags[len(frags)-1])

	return nil
}

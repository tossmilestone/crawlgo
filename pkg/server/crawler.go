package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/tossmilestone/crawlgo/pkg/util"
	"github.com/tossmilestone/crawlgo/pkg/web"
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
	webRender  web.Render
	total      int
	failed     int
	downloaded int
	progress   chan downloadTask
}

type downloadTask struct {
	urlPath string
	success bool
}

// DefaultSaveDir defines the default save directory for downloaded files.
var DefaultSaveDir = "./crawlgo"

// NewCrawler creates a Crawler object.
func NewCrawler(config *Config) (*Crawler, error) {
	// Create the save dir
	if config.SaveDir == "" {
		config.SaveDir = DefaultSaveDir
	}
	err := util.MkdirAll(config.SaveDir, os.ModeDir)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Init workers
	if config.Workers == 0 {
		config.Workers = runtime.NumCPU()
	}

	return &Crawler{
		config:     config,
		webRender:  web.NewPhantom(),
		total:      0,
		downloaded: 0,
		progress:   make(chan downloadTask),
	}, nil
}

// Run runs a crawler to execute the crawl tasks.
func (c *Crawler) Run(stop chan struct{}) {
	log.Print("Running crawler...")

	go func() {
		c.startCrawl(stop)
	}()

	go func() {
		c.displayProgress(stop)
	}()

	defer c.webRender.Stop()
	<-stop

	log.Print("Crawler stopped")
}

func (c *Crawler) startCrawl(stop chan struct{}) error {
	// Site analysis
	c.webRender.Run()
	log.Printf("Opening %s ...", c.config.Site)
	links, err := c.webRender.ExtractLinksFromSelector(c.config.Site, c.config.DownloadSelector)
	if err != nil {
		return err
	}
	c.total = len(links)
	log.Printf("Downloading %d links", c.total)
	if c.total == 0 {
		stop <- struct{}{}
		return nil
	}
	c.parallelizeDownload(links)
	return nil
}

func (c *Crawler) displayProgress(stop chan struct{}) {
	var failedDownload []downloadTask
	for {
		select {
		case result := <-c.progress:
			if !result.success {
				c.failed++
				failedDownload = append(failedDownload, result)
			} else {
				c.downloaded++
			}
			log.Printf("Downloaded %d, failed %d links of total %d", c.downloaded, c.failed, c.total)
			if c.downloaded + c.failed == c.total {
				if len(failedDownload) > 0 {
					log.Printf("Download failed of: %s", failedDownload)
				}
				log.Printf("All done, exit.")
				stop <- struct{}{}
			}
		case <-stop:
			return
		}
	}
}

func (c *Crawler) parallelizeDownload(downloads []interface{}) {
	downloadFile := func(index int) {
		download := downloadTask{
			urlPath: downloads[index].(string),
			success: true,
		}
		log.Printf("Start to download %s", downloads[index])
		err := c.download(downloads[index].(string))
		if err != nil {
			log.Print(err)
			download.success = false
			c.progress <- download
		} else {
			c.progress <- download
		}
	}
	util.Parallelize(c.config.Workers, len(downloads), downloadFile)
}

func (c *Crawler) download(urlStr string) error {
	frags := strings.Split(urlStr, "/")
	fileName := c.config.SaveDir + "/" + frags[len(frags)-1]

	file, err := util.Stat(fileName)
	if file != nil {
		log.Printf("%s exist", fileName)
		return nil
	}

	httpClient := &http.Client{
		Timeout: 0,
	}

	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(`download failed of "%s"`, frags[len(frags)-1])
	}

	out, err := util.Create(fileName)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Printf(`Downloaded of "%s"`, frags[len(frags)-1])

	return nil
}

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
)

type Config struct {
	SaveDir          string
	Site             string
	Workers          int
	DownloadSelector string
}

type Crawler struct {
	config         *Config
	phantom        *js.Phantom
	total          int
	downloaded     int
	progress       chan int
}

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
		config:         config,
		phantom:        js.NewPhantom(),
		total:          0,
		downloaded:     0,
		progress:       make(chan int),
	}, nil
}

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
		case <-c.progress:
			c.downloaded++
			log.Printf("Downloaded %d links of %d", c.downloaded, c.total)
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
		}
		c.progress <- 0
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

	out, err := os.Create(fileName)

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Printf(`Downloaded of "%s"`, frags[len(frags)-1])

	return nil
}

package server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tossmilestone/crawlgo/pkg/util"
	"github.com/tossmilestone/crawlgo/pkg/web"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"testing"
)

func mkdirAllOK(path string, perm os.FileMode) error {
	return nil
}

func mkdirAllErr(path string, perm os.FileMode) error {
	return fmt.Errorf("execute MkdirAll failed")
}

func statNotExist(name string) (os.FileInfo, error) {
	return nil, fmt.Errorf("file not exist")
}

func createOk(name string) (io.Writer, error) {
	return bytes.NewBuffer(make([]byte, 5)), nil
}

type fakeWebRender struct {
	linksCount int
}

func (w *fakeWebRender) Run() error {
	return nil
}

func (w *fakeWebRender) Stop() {}

func (w *fakeWebRender) ExtractLinksFromSelector(pageURL string, selector string) ([]interface{}, error) {
	var links []interface{}
	for i := 0; i < w.linksCount; i++ {
		links = append(links, "http://127.0.0.1:8080/test")
	}
	return links, nil
}

type fakeHTTPServer struct {
	server *http.Server
}

func (h *fakeHTTPServer) handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "test")
}

func (h *fakeHTTPServer) listen() {
	router := mux.NewRouter()
	router.HandleFunc("/test", h.handleTest).Methods("GET")
	h.server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Print("Test http server listened on :8080 ...")
	h.server.ListenAndServe()
}

func TestNewCrawler(t *testing.T) {
	tests := []struct {
		name       string
		saveDir    string
		mkdirAllOK bool
		workers    int
	}{
		{
			name:       "test 1",
			saveDir:    "",
			mkdirAllOK: false,
			workers:    0,
		},
		{
			name:       "test 2",
			saveDir:    "./test",
			mkdirAllOK: true,
			workers:    3,
		},
		{
			name:       "test 3",
			saveDir:    "",
			mkdirAllOK: true,
			workers:    0,
		},
	}

	for _, test := range tests {
		if test.mkdirAllOK {
			util.MkdirAll = mkdirAllOK
		} else {
			util.MkdirAll = mkdirAllErr
		}
		cfg := &Config{
			SaveDir: test.saveDir,
			Workers: test.workers,
		}
		c, err := NewCrawler(cfg)
		if err != nil {
			if test.mkdirAllOK {
				t.Errorf("Failed %s, expected not nil Crawler but nil", test.name)
			}
		} else {
			if !test.mkdirAllOK {
				t.Errorf("Failed %s, expected error", test.name)
			}

			expectedSaveDir := test.saveDir
			if test.saveDir == "" {
				expectedSaveDir = DefaultSaveDir
			}
			if c.config.SaveDir != expectedSaveDir {
				t.Errorf("Failed %s, expected save directory \"%s\", got \"%s\"",
					test.name, expectedSaveDir, c.config.SaveDir)
			}

			expectedWorkers := test.workers
			if test.workers == 0 {
				expectedWorkers = runtime.NumCPU()
			}
			if c.config.Workers != expectedWorkers {
				t.Errorf("Failed %s, expected %d workers, got %d",
					test.name, expectedWorkers, c.config.Workers)
			}

			if _, ok := c.webRender.(*web.Phantom); !ok {
				t.Errorf("Failed %s, expect that webRender is a web.Phantom",
					test.name)
			}
		}
	}
	// Recover the MkdirAll function back lest affecting other tests
	util.MkdirAll = util.DefaultMkdirAll
}

func TestCrawlerRun(t *testing.T) {
	tests := []struct {
		name  string
		links int
	}{
		{
			name:  "test 1",
			links: 0,
		},
		{
			name:  "test 2",
			links: 3,
		},
	}

	fakeH := &fakeHTTPServer{}
	go func() {
		fakeH.listen()
	}()

	for _, test := range tests {
		util.MkdirAll = mkdirAllOK
		util.Stat = statNotExist
		util.Create = createOk
		c, err := NewCrawler(&Config{
			SaveDir:          "",
			Site:             "http://127.0.0.1:8080/",
			Workers:          3,
			DownloadSelector: ".test",
		})
		if err != nil {
			t.Errorf("Failed %s, expected not nil Crawler", test.name)
		} else {
			c.webRender = &fakeWebRender{linksCount: test.links}
			stop := make(chan struct{})
			c.Run(stop)
			if test.links != c.downloaded {
				t.Errorf("Failed %s, expected downloaed %d, got %d",
					test.name, test.links, c.downloaded)
			}
		}
	}
	util.MkdirAll = util.DefaultMkdirAll
	util.Stat = util.DefaultStatFunc
	util.Create = util.DefaultCreateFunc
}

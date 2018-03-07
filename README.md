# Crawlgo [![CircleCI](https://circleci.com/gh/tossmilestone/crawlgo.svg?style=shield)](https://circleci.com/gh/tossmilestone/crawlgo) [![Coverage Status](https://coveralls.io/repos/github/tossmilestone/crawlgo/badge.svg?branch=master)](https://coveralls.io/github/tossmilestone/crawlgo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/tossmilestone/crawlgo)](https://goreportcard.com/report/github.com/tossmilestone/crawlgo)

Crawlgo is a crawler written in golang, it aims to be an extensible, scalable and high-performance distributed crawler system.

Using `phantomjs`, crawlgo can crawl web pages rendered with javascript.

## Prerequisite

* Linux OS
* phantomjs: `phantomjs` should be able to run through the env `PATH`. It can be downloaded [here](http://phantomjs.org/download.html).

## Install

```
go get github.com/tossmilestone/crawlgo
cd ${GOPATH}/src/github.com/tossmilestone/crawlgo
sudo make install
```

The above commands will install `crawlgo` in `${GOPATH}/go/bin`.

## Usage

```
crawlgo [flags]

Flags:
      --download-selector string   The DOM selector to query the links that will be downloaded from the site
      --enable-profile             enable profiling the program to start a pprof HTTP server on localhost:6360
  -h, --help                       help for crawlgo
      --save-dir string            The directory to save downloaded files. (default "./crawlgo")
      --site string                The site to crawl
      --version                    version for crawlgo
      --workers int                The number of workers to run the crawl tasks. If no set, will be 'runtime.NumCPU()'
```

Crawlgo uses file name to identify the downloaded links. If the file of a link is existed in the save directory, the link will be assumed downloaded already.
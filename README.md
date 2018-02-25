# Crawlgo [![CircleCI](https://circleci.com/gh/tossmilestone/crawlgo.svg?style=shield)](https://circleci.com/gh/tossmilestone/crawlgo)

Crawlgo is a crawler written in golang, it aims to be an extensible, scalable and high-performance distributed crawler system.


## Prerequisite

* Linux OS
* phantomjs: `phantomjs` binary must be in the `PATH` env

## Install

```
go get github.com/tossmilestone/crawlgo
```

## Usage

```
crawlgo --site {site_url} --workers {workers} --download-selector {download_selector} --save-dir {save_dir} --version
```


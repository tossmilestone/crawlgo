BIN := crawlgo

ARCH := amd64

PKG := github.com/tossmilestone/crawlgo

PACKAGES := $(shell go list ./... | grep -v /vendor/)

# Race detector is only supported on amd64.
RACE := $(shell test $$(go env GOARCH) != "amd64" || (echo "-race"))

VERSION := 1.0.0

SRC_DIRS := cmd pkg

INSTALL := n

.PHONY: setup ci check lint

all: setup install ci

ci: setup build check test coverage

install: INSTALL=y
install: build

build: bin/$(BIN)
	
setup:
	@go get -u github.com/golang/lint/golint
	@go get -u github.com/haya14busa/goverage
	@go get github.com/mattn/goveralls

bin/$(BIN):
	@echo "building $@"
	ARCH=$(ARCH)       \
	VERSION=$(VERSION) \
	PKG=$(PKG)         \
	INSTALL=$(INSTALL) \
	./build/build.sh

check: lint

lint:
	@test -z "$$(golint ./... | grep -v vendor/ | tee /dev/stderr)"

test:
	@go test -parallel 8 ${RACE} ${PACKAGES}

coverage:
	goverage -v -covermode=count -coverprofile=coverage.out ${PACKAGES}
	goveralls -coverprofile=coverage.out -service=circle-ci -repotoken ${COVERALLS_TOKEN}

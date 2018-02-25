BIN := crawlgo

ARCH := amd64

PKG := github.com/tossmilestone/crawlgo

VERSION := 1.0.0

SRC_DIRS := cmd pkg

INSTALL := n

.PHONY: setup ci check lint

all: setup install ci

ci: setup build check

install: INSTALL=y
install: build

build: bin/$(BIN)
	
setup:
	@go get -u github.com/golang/lint/golint

bin/$(BIN):
	@echo "building $@"
	ARCH=$(ARCH)       \
	VERSION=$(VERSION) \
	PKG=$(PKG)         \
	INSTALL=$(INSTALL) \
	./build/build.sh

check: lint

lint:
	@golint $(GOPACKAGES)

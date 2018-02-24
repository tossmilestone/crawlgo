BIN := crawlgo

ARCH := amd64

PKG := github.com/tossmilestone/crawlgo

VERSION := 1.0.0

SRC_DIRS := cmd pkg

all: build

build: bin/$BIN

bin/$BIN:
	@echo "building $@"
	ARCH=$(ARCH)       \
	VERSION=$(VERSION) \
	PKG=$(PKG)         \
	./build/build.sh

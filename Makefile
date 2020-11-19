FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
NO_C = CGO_ENABLED=0
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin
FLAGS_WIN = $(FLAGS_SHARED) GOOS=windows

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
FLAGS_LD=-ldflags "-X github.com/gnames/levenshtein.Version=${VERSION}"

GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GOCLEAN=$(GOCMD) clean
GOGET = $(GOCMD) get

all: install

test: deps install
	$(FLAG_MODULE) go test ./... -v

deps:
	$(GOCMD) mod download;

build:
	cd fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(NO_C) $(GOBUILD);

release:
	cd fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(NO_C) $(GOBUILD); \
	tar zcvf /tmp/fzdiff-${VER}-linux.tar.gz fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(NO_C) $(GOBUILD); \
	tar zcvf /tmp/fzdiff-${VER}-mac.tar.gz fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(NO_C) $(GOBUILD); \
	zip -9 /tmp/fzdiff-${VER}-win.zip fzdiff.exe; \
	$(GOCLEAN);

install:
	cd fzdiff; \
	$(FLAGS_SHARED) $(NO_C) $(GOINSTALL);


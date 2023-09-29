PROJ_NAME = levenshtein

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

NO_C = CGO_ENABLED=0
FLAGS_SHARED = $(NO_C) GOARCH=amd64
FLAGS_LINUX = $(FLAGS_SHARED) GOOS=linux
FLAGS_MAC = $(FLAGS_SHARED) GOOS=darwin
FLAGS_WIN = $(FLAGS_SHARED) GOOS=windows
FLAGS_LD=-ldflags "-w -s \
                  -X github.com/gnames/$(PROJ_NAME).Build=$(DATE) \
                  -X github.com/gnames/$(PROJ_NAME).Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w \
						-X github.com/gnames/$(PROJ_NAME).Build=$(DATE)"


GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOCLEAN=$(GOCMD) clean
GOGET = $(GOCMD) get

all: install

test: deps install
	@echo Run tests
	$(GOCMD) test -shuffle=on -count=1 -race -coverprofile=coverage.txt ./...

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download;

build:
	cd fzdiff; \
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD);

buildrel:
	cd fzdiff; \
	$(GOCLEAN); \
	$(NO_C) $(GORELEASE);

release:
	cd fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_LINUX) $(GOBUILD); \
	tar zcvf /tmp/fzdiff-${VER}-linux.tar.gz fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_MAC) $(GOBUILD); \
	tar zcvf /tmp/fzdiff-${VER}-mac.tar.gz fzdiff; \
	$(GOCLEAN); \
	$(FLAGS_WIN) $(GOBUILD); \
	zip -9 /tmp/fzdiff-${VER}-win.zip fzdiff.exe; \
	$(GOCLEAN);

install:
	cd fzdiff; \
	$(NO_C) $(GOINSTALL);


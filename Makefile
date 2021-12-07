.DEFAULT_GOAL := build

.PHONY: clean build docker-image test git-tag

BINARY        ?= ybdb-go-cli
SOURCES        = $(shell find . -name '*.go' | grep -v /vendor/)
VERSION       ?= $(shell git describe --tags --always --dirty)
GOPKGS         = $(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS   ?=
LDFLAGS       ?= -X github.com/radekg/yugabyte-db-go-client/config.Version=$(VERSION) -w -s
GOARCH        ?= amd64
GOOS          ?= linux

DOCKER_IMAGE_REPO ?= local/
CURRENT_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
TAG_VERSION ?= $(shell cat $(CURRENT_DIR)/.version)

TEST_TIMEOUT ?=120s

default: build

test:
	go clean -testcache
	go test -timeout ${TEST_TIMEOUT} -cover -v ./...

build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

docker-image:
	docker build -t $(DOCKER_IMAGE_REPO)$(BINARY):${TAG_VERSION} .

clean:
	@rm -rf build

git-tag:
	git tag v$(TAG_VERSION)
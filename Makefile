.DEFAULT_GOAL := build

.PHONY: clean build docker-image

BINARY        ?= ybdb-go-cli
SOURCES        = $(shell find . -name '*.go' | grep -v /vendor/)
VERSION       ?= $(shell git describe --tags --always --dirty)
GOPKGS         = $(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS   ?=
LDFLAGS       ?= -X github.com/radekg/yugabyte-db-go-client/config.Version=$(VERSION) -w -s
GOARCH        ?= amd64
GOOS          ?= linux

DOCKER_IMAGE_REPO ?= local/
DOCKER_IMAGE_VERSION ?= 0.0.1

default: build

test:
	go test -v -count=1 ./...

build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

docker-image:
	docker build -t $(DOCKER_IMAGE_REPO)$(BINARY):${DOCKER_IMAGE_VERSION} .

clean:
	@rm -rf build
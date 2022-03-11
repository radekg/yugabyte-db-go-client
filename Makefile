.DEFAULT_GOAL := test

.PHONY: clean test git-tag

CURRENT_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
TAG_VERSION ?= $(shell cat $(CURRENT_DIR)/.version | head -n1)

TEST_TIMEOUT ?=120s

default: test

test:
	go clean -testcache
	go test -timeout ${TEST_TIMEOUT} -cover -v ./...

clean:
	@rm -rf build

git-tag:
	git tag v$(TAG_VERSION)
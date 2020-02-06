SHELL=/bin/sh
IMAGE_TAG := $(shell git rev-parse HEAD)
export GO111MODULE=on

ifneq ($(version),)
#if version is set - tag image with given version
	IMAGE_TAG := $(version)
endif

.PHONY: test
test:
	go test -v -mod=vendor -cover -count=1 ./...

.PHONY: lint
lint:
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

.PHONY: deps
deps:
	rm -rf vendor
	go mod download
	go mod vendor
	go mod tidy

.PHONY: build
build:
	docker build  -t ivch/dynasty:latest  .

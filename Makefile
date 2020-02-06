SHELL=/bin/sh
IMAGE_TAG := $(shell git rev-parse HEAD)
SERVICES := `cat services`
export GO111MODULE=on

ifneq ($(version),)
#if version is set - tag image with given version
	IMAGE_TAG := $(version)
endif

#.PHONY: test
#test:
#	for i in ${SERVICES}; do \
#		go test -v -mod=vendor -cover -count=1 ./$$i; \
#	done
#
#.PHONY: lint
#lint:
#	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint
#	golangci-lint run

.PHONY: deps
deps:
	rm -rf vendor
	go mod download
	go mod vendor
	go mod tidy

#ifeq ($(service),)
#.PHONY: build
#build:
#	for i in ${SERVICES}; do \
#		docker build --build-arg SERVICE=$$i -t ivch/$$i:latest  . ; \
#	done
#else
#.PHONY: build
#build:
#	docker build --build-arg SERVICE=$(service) -t ivch/$(service):latest  .
#endif

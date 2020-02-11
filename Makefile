SHELL=/bin/sh
IMAGE_TAG := $(shell git rev-parse HEAD)
export GO111MODULE=on

ifneq ($(version),)
#if version is set - tag image with given version
	IMAGE_TAG := $(version)
endif

.PHONY: rundb
rundb:
	docker-compose -f docker-database.yml up -d

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
	docker build -t ivch/dynasty:latest  .

.PHONY: cover
cover:
	GO111MODULE=off go get github.com/axw/gocov/gocov
	GO111MODULE=off go get -u gopkg.in/matm/v1/gocov-html
	${GOPATH}/bin/gocov test ./... | ${GOPATH}/bin/gocov-html > coverage.html
	open coverage.html

.PHONY: gen
gen:
	GO111MODULE=off go get github.com/matryer/moq
	${GOPATH}/bin/moq -out modules/users/mock_test.go modules/users userRepository Service
	${GOPATH}/bin/moq -out modules/auth/mock_test.go modules/auth userService authRepository Service
	${GOPATH}/bin/moq -out modules/requests/mock_test.go modules/requests requestsRepository Service
	${GOPATH}/bin/moq -out clients/users/mock_test.go clients/users userService
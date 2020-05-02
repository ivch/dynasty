SHELL=/bin/sh
IMAGE_TAG := $(shell git rev-parse --short HEAD)
IMAGE_NAME = ivch/dynasty
export GO111MODULE=on

ifneq ($(version),)
#if version is set - tag image with given version
	IMAGE_TAG := $(version)
endif

testtag:
	@echo ${IMAGE_TAG}

.PHONY: rundb
rundb:
	docker-compose -f docker-database.yml up -d --remove-orphans --force-recreate

.PHONY: test
test:
	go test -v -mod=vendor -cover -count=1 ./...

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

.PHONY: deps
deps:
	rm -rf vendor
	go mod download
	go mod vendor
	go mod tidy

.PHONY: build
build:
	tar cfz zoneinfo.tar.gz /usr/share/zoneinfo
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG}  .
	rm zoneinfo.tar.gz

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
	${GOPATH}/bin/moq -out modules/requests/mock_test.go modules/requests requestsRepository s3Client Service
	${GOPATH}/bin/moq -out modules/dictionaries/mock_test.go modules/dictionaries dictRepository Service
	${GOPATH}/bin/moq -out clients/users/mock_test.go clients/users userService

.PHONY: tag
tag:
	docker pull ${IMAGE_NAME}:latest
	docker tag ${IMAGE_NAME}:latest ${IMAGE_NAME}:prev
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest

.PHONY: push
push: tag
	docker push ${IMAGE_NAME}:prev
	docker push ${IMAGE_NAME}:${IMAGE_TAG}
	docker push ${IMAGE_NAME}:latest
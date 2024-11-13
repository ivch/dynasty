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
	golangci-lint run

.PHONY: inastall-lint
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: deps
deps:
	rm -rf vendor
	go mod download
	go mod tidy

.PHONY: build
build:
	tar cfz zoneinfo.tar.gz /usr/share/zoneinfo
	docker build --build-arg CODECOV_TOKEN=${DYN_CODECOV_TOKEN} -t ${IMAGE_NAME}:${IMAGE_TAG}  .
	rm zoneinfo.tar.gz

.PHONY: cover
cover:
	GO111MODULE=off go get github.com/axw/gocov/gocov
	GO111MODULE=off go get -u gopkg.in/matm/v1/gocov-html
	${GOPATH}/bin/gocov test ./... | ${GOPATH}/bin/gocov-html > coverage.html
	open coverage.html

.PHONY: gen
gen:
	go install github.com/matryer/moq@latest
	${GOPATH}/bin/moq -out server/handlers/users/mock_test.go server/handlers/users userRepository mailSender
	${GOPATH}/bin/moq -out server/handlers/users/transport/mock_test.go server/handlers/users/transport UsersService
	${GOPATH}/bin/moq -out common/clients/users/mock_test.go common/clients/users userService
	${GOPATH}/bin/moq -out server/handlers/auth/transport/mock_test.go server/handlers/auth/transport AuthService
	${GOPATH}/bin/moq -out server/handlers/auth/mock_test.go server/handlers/auth userService authRepository
	${GOPATH}/bin/moq -out server/handlers/dictionaries/mock_test.go server/handlers/dictionaries dictRepository
	${GOPATH}/bin/moq -out server/handlers/dictionaries/transport/mock_test.go server/handlers/dictionaries/transport DictionaryService
	${GOPATH}/bin/moq -out server/handlers/requests/transport/mock_test.go server/handlers/requests/transport RequestsService
	${GOPATH}/bin/moq -out server/handlers/requests/mock_test.go server/handlers/requests requestsRepository s3Client

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

BINARY_NAME=eetest

REPO_NAME=eetestproject
REPO=registry.foo.com/foo/repos/rp/eetestproject
BUILD_TIME=$(shell date -u +"%FT%TZ")
ifndef HASH_FROM_DOCKER
GIT_HASH=$(shell git rev-parse HEAD)
else
GIT_HASH=${HASH_FROM_DOCKER}
endif


.PHONY: run
run:
	go run \
		-mod=vendor \
		-ldflags="-s -w -X main.gitHash=${GIT_HASH} -X main.buildTime=${BUILD_TIME}" \
		cmd/$(BINARY_NAME)/main.go

.PHONY: linux
linux:
	CGO_ENABLED=0 \
	GOOS=linux \
	GO111MODULE=auto \
	GOBIN=${CURDIR} \
	go install \
		-mod=vendor \
		-ldflags="-s -w -X main.gitHash=${GIT_HASH} -X main.buildTime=${BUILD_TIME}" \
		-trimpath \
		./cmd/...

.PHONY: test
test:
	go test ./... -cover -coverprofile coverage.out -race -mod=vendor -v

.PHONY: build-docker
build-docker:
	docker build --build-arg GIT_HASH=${GIT_HASH} -t ${REPO} .

.PHONY: lint
lint:
	docker run \
		--rm \
		-it \
		-w /go/src/repos/RP/${REPO_NAME}.git \
		-v $(shell pwd):/go/src/repos/RP/${REPO_NAME}.git \
		golangci/golangci-lint:v1.22 \
		golangci-lint run \
			-v \
			-c .golangci.yml

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: gofmt
gofmt:
	gofmt -l -s -w ./pkg ./cmdTt
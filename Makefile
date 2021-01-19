NAME := rabbitmq-publisher
VERSION := 0.0.1
REVISION := $(shell git rev-parse --short HEAD)
BRANCH := $(shell git symbolic-ref --short HEAD)
DATE            := $(shell date -uR)
GOVERSION       := $(shell go version)
GITHUB_USER     := ken-aio
OSARCH          := "darwin/amd64 linux/amd64"

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := "-X \"github.com/ken-aio/rabbitmq-publisher/cmd.Revision=${REVISION}\" -X \"github.com/ken-aio/rabbitmq-publisher/cmd.BuildDate=${DATE}\" -X \"github.com/ken-aio/rabbitmq-publisher/cmd.GoVersion=${GOVERSION}\" -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

$(GOX): ; @go get github.com/mitchellh/gox
$(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/arc
$(GHR): ; @go get github.com/tcnksm/ghr

GHR := $(GOBIN)/ghr

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show help see: https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build app for developers os
	GO111MODULE=on CGO_ENABLED=0 go build -mod=readonly $(LDFLAGS) -o dist/$(NAME)

.PHONY: cross-build
cross-build: $(GOX) ## build some arch
	rm -rf ./out && \
	gox -ldflags $(LDFLAGS) -osarch $(OSARCH) -output "./out/${NAME}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: package
package: cross-build $(ARCHIVER) ## package bin
	rm -rf ./pkg && mkdir ./pkg && \
	pushd out && \
	find * -type d -exec arc archive ../pkg/{}.tar.gz {}/$(NAME) \; && \
	popd

.PHONY: release
release: $(GHR) ## release to github
	ghr -u $(GITHUB_USER) $(VERSION) pkg/

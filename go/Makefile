.DEFAULT_GOAL := build

export CGO_ENABLED=0
export VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
export SUFFIX=$(shell git describe --exact-match > /dev/null 2>&1 || echo "-dev")
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
export BUILDDATE=$(shell LANG=en_us_88591 date)
BINARY={{ .name }}
LDFLAGS=-ldflags "-X '{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd.Version=$(VERSION)$(SUFFIX)' \
		-X '{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd.Build=$(BUILD)' \
		-X '{{ .gitserver }}/{{ .owner }}/{{ .name }}/cmd.BuildDate=$(BUILDDATE)' -s -w"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build
	go build $(LDFLAGS) -o $(BINARY)

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: tmp
tmp: ## Build and output the binary in /tmp
	go build $(LDFLAGS) -o /tmp/$(BINARY)
{{ if .docker }}
.PHONY: docker
docker: ## Build the docker image
	docker build -t $(BINARY):latest $(if $(filter undefined,$(BUILD)),,-t $(BINARY):$(BUILD) )-f Dockerfile .
{{ end }}{{ if .goreleaser }}
.PHONY: release
release: ## Create a new release on Github
	VERSION="$(VERSION)" BUILD="$(BUILD)" BUILDDATE="$(BUILDDATE)" goreleaser release

.PHONY: snapshot
snapshot: ## Create a new snapshot release
	VERSION="$(VERSION)" BUILD="$(BUILD)" BUILDDATE="$(BUILDDATE)" goreleaser release --snapshot --clean
{{ end }}
.PHONY: lint
lint: ## Runs the linter
	$(GOPATH)/bin/golangci-lint run --exclude-use-default=false

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test -race -coverprofile="coverage.txt" ./...

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi

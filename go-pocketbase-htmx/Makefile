.DEFAULT_GOAL := build

export CGO_ENABLED=0
export VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
export SUFFIX=$(shell git describe --exact-match > /dev/null 2>&1 || echo "-dev")
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
export BUILDDATE=$(shell LANG=en_us_88591 date)
BINARY={{.name}}
LDFLAGS=-ldflags "-X '{{.gitserver}}/{{.owner}}/{{.name}}/cmd.Version=$(VERSION)$(SUFFIX)' \
		-X '{{.gitserver}}/{{.owner}}/{{.name}}/cmd.Build=$(BUILD)' \
		-X '{{.gitserver}}/{{.owner}}/{{.name}}/cmd.BuildDate=$(BUILDDATE)' -s -w"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: bootstrap
bootstrap: ## Install TemplUI components, download HTMX, and tidy dependencies
	templui init
	templui add sidebar card button input form label switch toast icon avatar dropdown
	curl -sL https://cdn.jsdelivr.net/npm/htmx.org@2.0.10/dist/htmx.min.js -o assets/js/htmx.min.js
	templ generate ./...
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify
	go mod tidy

.PHONY: assets
assets: ## Build assets
	templ generate
	tailwindcss -i ./assets/css/input.css -o ./assets/css/output.css --minify

.PHONY: build
build: assets ## Build binary
	go build $(LDFLAGS) -o $(BINARY)

.PHONY: run
run: assets ## Run in development mode
	go run . serve --dev

.PHONY: docker
docker: ## Build the docker image
	docker build -t $(BINARY):latest $(if $(filter undefined,$(BUILD)),,-t $(BINARY):$(BUILD) )-f Dockerfile .

.PHONY: release
release: ## Create a new release on Github
	goreleaser release

.PHONY: snapshot
snapshot: ## Create a new snapshot release
	goreleaser release --snapshot --clean

.PHONY: lint
lint: ## Runs the linter
	golangci-lint run

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test -race -coverprofile="coverage.txt" ./...

.PHONY: clean
clean: ## Remove the binary and data
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi
	if [ -f coverage.txt ] ; then rm coverage.txt ; fi

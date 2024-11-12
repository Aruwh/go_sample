CMD	 	 = ./cmd
MODULES  = $(shell find $(CMD) -mindepth 1 -maxdepth 1 -type d)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell cat VERSION || NOT_AVAILABLE)
BRANCH ?= $(shell git branch --show-current 2> /dev/null || echo 0)
LAST_COMMIT=$(shell git rev-parse --short HEAD 2> /dev/null || echo 0)
SERVICE_VERSION = $(BRANCH)-$(VERSION).$(LAST_COMMIT)
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./... | grep -v /mock))
BIN		 = $(CURDIR)/bin
GO       = go

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

GO111MODULE=on

.PHONY: all
all: fmt lint vendor | $(BIN) ; $(info $(M) building osx releases...) @ ## Build osx release binaries
	$Q $(foreach MODULE,$(MODULES), \
		$(info $(M) building $(shell basename $(MODULE)))	\
		$(GO) build \
			-ldflags '-X main.Version=$(SERVICE_VERSION) -X main.BuildDate=$(DATE) -s -w' \
			-tags release \
			-o $(BIN)/$(shell basename $(MODULE)) $(MODULE);)

.PHONY: release
release: fmt lint vendor | $(BIN) ; $(info $(M) building linux release...) @ ## Build linux release binary
	$Q $(foreach MODULE,$(MODULES), \
    	$(info $(M) building $(shell basename $(MODULE))) \
    	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build \
			-a -installsuffix cgo \
        	-ldflags '-X main.Version=$(SERVICE_VERSION) -X main.BuildDate=$(DATE) -s -w -extldflags "-static"'	\
        	-tags release	\
        	-o $(BIN)/$(shell basename $(MODULE)) $(MODULE);)

# Tools

$(BIN):
	@mkdir -p $@


$(GOBIN)/%: ; $(info $(M) installing $(PACKAGE)...)
	$Q $(GO) install $(ARGS) $(PACKAGE)

$(BIN)/%: | $(BIN) ; $(info $(M) building $(PACKAGE)...)
	$Q tmp=$$(mktemp -d); \
    	env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
        || ret=$$?; \
        rm -rf $$tmp ; exit $$ret

GOIMPORTS = $(BIN)/goimports
$(BIN)/goimports: PACKAGE=golang.org/x/tools/cmd/goimports

GOLINT = $(GOBIN)/golangci-lint
$(GOBIN)/golangci-lint: PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.1

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go-junit-report
$(BIN)/go-junit-report: PACKAGE=github.com/jstemmer/go-junit-report

MOCKERY = $(GOBIN)/mockery
$(GOBIN)/mockery: PACKAGE=github.com/vektra/mockery/v2@v2.29.0

# Dependencies / Tools
COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html

.PHONY: test tests
test tests: fmt ; $(info $(M) running unit tests...) @ ## Run unit tests
	$Q $(GO) test -tags development -count=1 $(PKGS)

.PHONY: test-cover test-coverage test-coverage-tools
test-coverage-tools: $(GOCOV) $(GOCOVXML) $(GO2XUNIT)
test-cover test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-cover test-coverage: fmt test-coverage-tools ; $(info $(M) running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test -v -tags=release,integration \
		-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $(PKGS) | \
					grep '^$(MODULE)/' | \
					tr '\n' ',' | sed 's/,$$//') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(PKGS)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
	## Open coverage.out with: $ go tool cover -html=$(COVERAGE_PROFILE)

.PHONY: test-pipeline
test-pipeline: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-pipeline: fmt test-coverage-tools ; $(info $(M) running tests and generate code coverage…) @ ## Run tests & generate code coverage for pipeline
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test -v -tags=release,integration \
		-coverpkg=$$($(GO) list -f '{{ join .Deps "\n" }}' $(PKGS) | \
					grep '^$(MODULE)/' | \
					tr '\n' ',' | sed 's/,$$//') \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(PKGS) | tee test/tests.output
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)
	$Q cat test/tests.output | $(GO2XUNIT) > test/tests.xml

.PHONY: init
init: deps gen vendor ; $(info $(M) initialising the project...) @ ## Initialising the project

.PHONY: lint lint-staged lint-no-fix
lint-no-fix: $(GOLINT) ; $(info $(M) running golangci-lint (without fixing issues)…) @ ## Run golangci-lint on codebase (without fixing issues)
	$Q $(GOLINT) run --verbose --sort-results
lint-staged: NEW=--new ; $(info $(M) running golangci-lint on staged files...) ## Run golangci-lint on staged source files only
lint-staged: lint
lint: $(GOLINT) ; $(info $(M) running golangci-lint…) @ ## Run golangci-lint on codebase
	$Q $(GOLINT) run --verbose --sort-results --fix $(NEW)

.PHONY: fmt
fmt: $(GOIMPORTS) ; $(info $(M) running gofmt…) @ ## Run gofmt & goimports on all source files
	$Q $(GO) fmt $(PKGS)
	$Q $(GOIMPORTS) -w .

.PHONY: vendor
vendor: ; $(info $(M) vendoring dependencies...) @ ## Tidying modules & setting up go dependencies
	$Q $(GO) mod tidy
	$Q $(GO) mod vendor

.PHONY: install deps
deps: ; $(info $(M) installing dependencies & tooling...) @ ## Installing dependencies & tooling
	$Q ./scripts/install-deps.sh

.PHONY: generate gen
gen: ; $(info $(M) generating swagger code and source files...)	@ ## Run all source code generation tools
	$Q ./scripts/generate.sh
	$Q rm -f internal/mocks/* && $(GO) generate ./...

.PHONY: clean all
clean: ; $(info $(M) cleaning artifacts...)	@ ## Cleanup everything
	$Q rm -rf $(BIN)

.PHONY: help
help:
	$Q grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: package
package: ; $(info $(M) building & pushing docker image...)	@ ## Build & publish docker image
	$Q ./scripts/package.sh

.PHONY: up
up: ; $(info $(M) starting local docker-compose cluster...) @ ## Starts all services
	$Q docker network create fewoserv || true
	$Q docker-compose -f deploy/docker-compose.yaml up --remove-orphans -d

.PHONY: down
down: ; $(info $(M) stoping local docker-compose cluster...) @ ## Stops all services
	$Q docker-compose -f deploy/docker-compose.yaml down

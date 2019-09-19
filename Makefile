# Default build target
.DEFAULT_GOAL = default

# Package variables
START_COMMAND = main.go start
MODULE        = $(shell env $(GO) list -m)
DATE          = $(shell date +%FT%T%z)
VERSION       = $(shell cat VERSION)
BIN           = $(CURDIR)/bin
GO            = go

# Helper variables
V             = 0
Q             = $(if $(filter 1,$V),,@)
M             = $(shell printf "\033[34;1m▶\033[0m")
PKGS          = $(shell $(GO) list ./...)
TESTPKGS      = $(shell $(GO) list -f \
			    '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			    $(PKGS))
TIMEOUT       = 15

# Run

start: ; $(info $(M) Starting app...) @ ## Start the app
	$Q $(GO) run $(START_COMMAND)

# Build
build: ; $(info $(M) Building image eye-of-sauron...) @ ## Build the docker image
	$Q docker build -t eye-of-sauron .

# Tools

$(BIN):
	@mkdir -p $@
$(BIN)/%: | $(BIN) ; $(info $(M) Building $(PACKAGE)…)
	$Q tmp=$$(mktemp -d); \
	env GO111MODULE=off GOPATH=$$tmp GOBIN=$(BIN) $(GO) get $(PACKAGE) \
	|| ret=$$?; \
	rm -rf $$tmp ; exit $$ret

GOLINT = $(BIN)/golint
$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOCOV = $(BIN)/gocov
$(BIN)/gocov: PACKAGE=github.com/axw/gocov/...

GOCOVXML = $(BIN)/gocov-xml
$(BIN)/gocov-xml: PACKAGE=github.com/AlekSi/gocov-xml

GO2XUNIT = $(BIN)/go2xunit
$(BIN)/go2xunit: PACKAGE=github.com/tebeka/go2xunit

# Tests

test: fmt lint ; $(info $(M) Running tests…) @ ## Run tests
	$Q $(GO) test -timeout $(TIMEOUT)s $(TESTPKGS)

TEST_TARGETS := test-default test-bench test-short test-verbose test-race
.PHONY: $(TEST_TARGETS) test-xml
test-bench:   ARGS=-run=__absolutelynothing__ -bench=. ## Run benchmarks
test-short:   ARGS=-short                              ## Run only short tests
test-verbose: ARGS=-v                                  ## Run tests in verbose mode with coverage reporting
test-race:    ARGS=-race                               ## Run tests with race detector
$(TEST_TARGETS): NAME=$(MAKECMDGOALS:test-%=%)
$(TEST_TARGETS): named-test
named-test: fmt lint ; $(info $(M) Running $(NAME:%=% )tests…)
	$Q $(GO) test -timeout $(TIMEOUT)s $(ARGS) $(TESTPKGS)

test-xml: fmt lint | $(GO2XUNIT) ; $(info $(M) Running xUnit tests…) @ Run tests with xUnit output
	$Q mkdir -p test
	$Q 2>&1 $(GO) test -timeout $(TIMEOUT)s -v $(TESTPKGS) | tee test/tests.output
	$(GO2XUNIT) -fail -input test/tests.output -output test/tests.xml

COVERAGE_MODE    = atomic
COVERAGE_PROFILE = $(COVERAGE_DIR)/profile.out
COVERAGE_XML     = $(COVERAGE_DIR)/coverage.xml
COVERAGE_HTML    = $(COVERAGE_DIR)/index.html
.PHONY: test-coverage test-coverage-tools
test-coverage-tools: | $(GOCOV) $(GOCOVXML)
test-coverage: COVERAGE_DIR := $(CURDIR)/test/coverage.$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
test-coverage: fmt lint test-coverage-tools ; $(info $(M) Running coverage tests…) @ ## Run coverage tests
	$Q mkdir -p $(COVERAGE_DIR)
	$Q $(GO) test \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile="$(COVERAGE_PROFILE)" $(TESTPKGS)
	$Q $(GO) tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	$Q $(GOCOV) convert $(COVERAGE_PROFILE) | $(GOCOVXML) > $(COVERAGE_XML)

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) Running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status

.PHONY: fmt
fmt: ; $(info $(M) Running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt

# Misc

default: help

.PHONY: clean
clean: ; $(info $(M) Cleaning…)	@ ## Cleanup everything
	@rm -rf $(BIN)
	@rm -rf test/tests.* test/coverage.*

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
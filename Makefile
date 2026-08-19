SHELL := /bin/bash
GO ?= go
GO_PACKAGES := ./...
UNIT_PACKAGES := ./internal/...
INTEGRATION_PACKAGES := ./test/integration/...
BIN_DIR := bin
COVERAGE_FILE ?= coverage.out
COVERAGE_THRESHOLD ?= 80.0
PACKAGE_FILE ?= $(BIN_DIR)/homework_go_5-linux-amd64.tar.gz
CMDS := 01_interfaces 02_methodsets 03_io 04_args 05_common

.PHONY: help deps-check mod-check fmt fmt-check vet test test-unit test-integration test-race coverage coverage-check build package clean run-all ci compile test-interfaces test-methodsets test-io test-args test-common $(addprefix run-,$(CMDS))

help:
	@echo "Available commands:"
	@echo "  make compile          - compile all packages without running tests"
	@echo "  make test-interfaces  - run interface tasks tests"
	@echo "  make test-methodsets  - run method set tasks tests"
	@echo "  make test-io          - run io tasks tests"
	@echo "  make test-args        - run command-line argument tasks tests"
	@echo "  make test-common      - run common file tasks tests"
	@echo "  make ci               - full local CI after solving all tasks"

deps-check:
	$(GO) mod download
	$(GO) mod verify

mod-check:
	$(GO) mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
		git diff --exit-code -- go.mod; \
		if [ -f go.sum ]; then git diff --exit-code -- go.sum; fi; \
	else \
		echo "Skipping git diff because this directory is not a git repository"; \
	fi

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './$(BIN_DIR)/*')

fmt-check:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './$(BIN_DIR)/*'))"; \
	if [ -n "$$files" ]; then echo "Go files are not formatted:"; echo "$$files"; exit 1; fi

vet:
	$(GO) vet $(GO_PACKAGES)

test: test-unit test-integration

test-unit:
	$(GO) test $(UNIT_PACKAGES)

test-integration:
	$(GO) test $(INTEGRATION_PACKAGES)

compile:
	$(GO) test -run '^$$' $(GO_PACKAGES)

test-interfaces:
	$(GO) test ./internal/interfaces/...

test-methodsets:
	$(GO) test ./internal/methodsets/...

test-io:
	$(GO) test ./internal/ioflow/...

test-args:
	$(GO) test ./internal/cliargs/...

test-common:
	$(GO) test ./internal/common/...

test-race:
	$(GO) test -race $(UNIT_PACKAGES)

coverage:
	$(GO) test $(UNIT_PACKAGES) -covermode=atomic -coverprofile=$(COVERAGE_FILE)
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check: coverage
	@coverage="$$(go tool cover -func=$(COVERAGE_FILE) | awk '/^total:/ {gsub("%", "", $$3); print $$3}')"; \
	awk -v coverage="$$coverage" -v threshold="$(COVERAGE_THRESHOLD)" 'BEGIN { \
		if (coverage + 0 < threshold + 0) { printf "coverage %.1f%% is below threshold %.1f%%\n", coverage, threshold; exit 1 } \
		printf "coverage %.1f%% is enough; threshold %.1f%%\n", coverage, threshold; \
	}'

run-all:
	@for cmd in $(CMDS); do \
		echo "== $$cmd =="; \
		if [ "$$cmd" = "04_args" ]; then $(GO) run ./cmd/$$cmd "Maria Petrova" 2; else $(GO) run ./cmd/$$cmd; fi; \
	done

run-01_interfaces:
	$(GO) run ./cmd/01_interfaces
run-02_methodsets:
	$(GO) run ./cmd/02_methodsets
run-03_io:
	$(GO) run ./cmd/03_io
run-04_args:
	$(GO) run ./cmd/04_args "Maria Petrova" 2
run-05_common:
	$(GO) run ./cmd/05_common

build:
	@mkdir -p $(BIN_DIR)
	@for cmd in $(CMDS); do $(GO) build -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; done

package: build
	@mkdir -p $(BIN_DIR)
	tar -czf $(PACKAGE_FILE) -C $(BIN_DIR) $(CMDS)

ci: deps-check mod-check fmt-check vet test-unit test-integration test-race coverage-check build package

clean:
	rm -rf $(BIN_DIR) $(COVERAGE_FILE)

GO_BUILD = go build
GO_TEST = go test
GO_TEST_FLAG = -v
GO_ENV = GO111MODULE=on
BINDIR = ./bin
CMDDIR = ./cmd
MAIN = main.go
BIN = llss
INTERNAL = llss/internal
TEST := $(INTERNAL)/scanner
TEST := $(TEST) $(INTERNAL)/analyzer
TEST := $(TEST) $(INTERNAL)/domainscanner
TEST_FILE := ./internal/scanner/scanner_test.go
TEST_FILE := $(TEST_FILE) ./internal/analyzer/analyzer_test.go
TEST_FILE := $(TEST_FILE) ./internal/domainscanner/domainscanner_test.go

.PHONY: all
all: build

.PHONY: build
build: $(CMDDIR)/$(BIN)/$(MAIN)
	$(GO_ENV) $(GO_BUILD) -o "$(BINDIR)/$(BIN)" "$(CMDDIR)/$(BIN)/$(MAIN)"

.PHONY: clean
clean:
	rm -rf $(BINDIR)

.PHONY: test
test: $(TEST_FILE)
	$(GO_ENV) $(GO_TEST) $(GO_TEST_FLAG) $(TEST)

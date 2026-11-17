APP_NAME     := imgcontent
DIST_DIR     := bin
OUTPUT       := $(DIST_DIR)/$(APP_NAME)
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "devel")
LDFLAGS      := -X main.Version=$(VERSION)
INSTALL_PATH ?= $(HOME)/bin

.PHONY: build
build: ## build binary
	go build -ldflags '$(LDFLAGS)' -o $(OUTPUT) .

.PHONY: test
test: ## run tests
	go test ./...

.PHONY: install
install: build ## install binary to INSTALL_PATH (default: ~/bin)
	install -m 755 $(OUTPUT) $(INSTALL_PATH)/$(APP_NAME)

.PHONY: clean
clean: ## remove build artifacts
	@test -n "$(DIST_DIR)" && rm -rf "$(DIST_DIR)"

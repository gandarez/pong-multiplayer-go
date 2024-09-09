.DEFAULT_GOAL := build

# globals
BINARY_NAME?=pongo
BUILD_DIR?="./build"
CGO_ENABLED?=1

# basic Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# get GOPATH, GOOS and GOARCH according to OS
ifeq ($(OS),Windows_NT) # is Windows_NT on XP, 2000, 7, Vista, 10...
    GOPATH=$(go env GOPATH)
	GOOS=$(shell cmd /c go env GOOS)
	GOARCH=$(shell cmd /c go env GOARCH)
else
    GOPATH=$(shell go env GOPATH)
	GOOS=$(shell go env GOOS)
	GOARCH=$(shell go env GOARCH)
endif

# targets
build-all: build-all-darwin build-all-linux build-all-windows

build-all-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) build

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) build

build-all-linux: build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-linux-riscv64

build-linux-386:
	GOOS=linux GOARCH=386 $(MAKE) build

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) build

build-linux-arm:
	GOOS=linux GOARCH=arm $(MAKE) build

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build

build-all-windows: build-windows-386 build-windows-amd64 build-windows-arm64

build-windows-386:
	GOOS=windows GOARCH=386 $(MAKE) build-windows

build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(MAKE) build-windows

build-windows-arm64:
	GOOS=windows GOARCH=arm64 $(MAKE) build-windows

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -v \
		-o ${BUILD_DIR}/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(FILE_EXT) ./cmd/game/main.go

build-windows:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -v \
		-o ${BUILD_DIR}/$(BINARY_NAME)-$(GOOS)-$(GOARCH).exe ./cmd/game/main.go

install: install-deps install-linter

.PHONY: install-linter
install-linter:
ifneq "$(INSTALLED_LINT_VERSION)" "$(LATEST_LINT_VERSION)"
	@echo "new golangci-lint version found:" $(LATEST_LINT_VERSION)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin latest
endif

.PHONY: install-deps
install-deps:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing Ebiten dependencies..."
	./scripts/install-dependencies.sh

.PHONY: test
test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

# run static analysis tools, configuration in ./.golangci.yml file
.PHONY: lint
lint: install-linter
	golangci-lint run ./...

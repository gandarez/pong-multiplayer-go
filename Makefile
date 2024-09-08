.PHONY: all build run install-deps test clean

BINARY_NAME=pongo

all: build

build:
	go build -o $(BINARY_NAME) cmd/game/main.go

run: build
	./$(BINARY_NAME)

install-deps:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing Ebiten dependencies..."
	bash scripts/install-dependencies.sh

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)


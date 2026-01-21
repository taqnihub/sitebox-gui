# SiteBox GUI Makefile

APP_NAME=sitebox-gui
VERSION?=1.0.0

.PHONY: all dev build build-all clean install frontend help

all: build

## Development

dev: ## Run in development mode with hot reload
	wails dev

## Build commands

build: ## Build for current platform
	wails build

build-darwin-arm64: ## Build for macOS Apple Silicon
	wails build -platform darwin/arm64

build-darwin-amd64: ## Build for macOS Intel
	wails build -platform darwin/amd64

build-windows: ## Build for Windows
	wails build -platform windows/amd64

build-linux: ## Build for Linux
	wails build -platform linux/amd64

build-all: ## Build for all platforms
	wails build -platform darwin/arm64
	wails build -platform darwin/amd64
	wails build -platform windows/amd64
	wails build -platform linux/amd64

## Frontend

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-build: ## Build frontend
	cd frontend && npm run build

frontend-dev: ## Run frontend dev server
	cd frontend && npm run dev

## Utility

clean: ## Clean build artifacts
	rm -rf build
	rm -rf frontend/dist
	rm -rf frontend/node_modules

generate: ## Generate Wails bindings
	wails generate module

tidy: ## Tidy Go modules
	go mod tidy

deps: ## Install all dependencies
	go mod download
	cd frontend && npm install

## Installation

install: build ## Build and install (macOS)
	@if [ -d "/Applications/SiteBox.app" ]; then \
		rm -rf "/Applications/SiteBox.app"; \
	fi
	@if [ -d "build/bin/SiteBox.app" ]; then \
		cp -r "build/bin/SiteBox.app" /Applications/; \
		echo "Installed SiteBox to /Applications"; \
	else \
		echo "Build not found. Run 'make build' first."; \
		exit 1; \
	fi

## Version

version: ## Show version
	@echo $(VERSION)

## Help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

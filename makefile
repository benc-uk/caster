SRC_DIR := main
GO_PKG := github.com/benc-uk/caster

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go

.PHONY: help run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 🔨 Build binaries
	@figlet $@
	@mkdir -p bin
	go mod tidy
	GOOS=linux go build -o bin/caster $(GO_PKG)/...
	GOOS=windows go build -o bin/caster.exe $(GO_PKG)/...

clean: ## ♻️  Clean up
	@figlet $@
	@rm -rf bin

lint: ## 🌟 Lint & format, will not fix but sets exit code on error
	@figlet $@ || true
	cd $(SRC_DIR); golangci-lint run --modules-download-mode=mod *.go


lint-fix: ## 🔍 Lint & format, will try to fix errors and modify code
	@figlet $@ || true

	cd $(SRC_DIR); golangci-lint run --modules-download-mode=mod *.go --fix

run: ## 🏃 Run application
	@figlet $@ || true
	air -c .air.toml

windows: ## 💻 Bundle Windows version
	@figlet $@
	make build
	cp bin/caster.exe /mnt/c/Temp/caster.exe
	cp -r ./textures /mnt/c/Temp/
	cp -r ./sprites /mnt/c/Temp/
	cp -r ./maps /mnt/c/Temp/
	cp -r ./sounds /mnt/c/Temp/
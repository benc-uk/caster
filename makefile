SRC_DIR := src
GO_PKG := github.com/benc-uk/caster
WIN_DIR := /mnt/c/Temp/caster
LINUX_DIR := /tmp/caster

# Things you don't want to change
REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
GOLINT_PATH := $(REPO_DIR)/bin/golangci-lint # Remove if not using Go

.PHONY: help run lint lint-fix
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build-linux: ## 🔨 Build binaries for Linux
	@figlet $@
	@mkdir -p bin
	go mod tidy
	GOOS=linux go build -o bin/caster $(GO_PKG)/...

build-win: ## 🔨 Build binaries for Windows
	@figlet $@
	@mkdir -p bin
	go mod tidy
	GOOS=windows go build -o bin/caster.exe $(GO_PKG)/...

build: build-win build-linux ## 🔨 Build binaries
	cp -r gfx editor/gfx

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

release-windows: build-win ## 💻 Bundle Windows version
	@figlet $@
	rm -rf $(WIN_DIR)/
	mkdir -p $(WIN_DIR)
	cp bin/caster.exe $(WIN_DIR)/caster.exe
	cp -r ./gfx $(WIN_DIR)/
	cp -r ./maps $(WIN_DIR)/
	cp -r ./sounds $(WIN_DIR)/
	cp -r ./fonts $(WIN_DIR)/
	cd $(WIN_DIR); zip ./crypt-caster-win.zip ./*

release-linux: build-linux ## 💻 Bundle Linux version
	@figlet $@
	rm -rf $(LINUX_DIR)/
	mkdir -p $(LINUX_DIR)
	cp bin/caster $(LINUX_DIR)/caster
	cp -r ./gfx $(LINUX_DIR)/
	cp -r ./maps $(LINUX_DIR)/
	cp -r ./sounds $(LINUX_DIR)/
	cp -r ./fonts $(LINUX_DIR)/
	cd $(LINUX_DIR); zip ./crypt-caster-linux.zip ./*
	cp $(LINUX_DIR)/crypt-caster-linux.zip $(WIN_DIR)/crypt-caster-linux.zip

run-editor: ## 📝 Run level editor
	@figlet $@
	cd editor; browser-sync start --server --single --no-ui --no-open --no-notify --watch
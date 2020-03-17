.DEFAULT_GOAL := help

TAG      = $(shell git describe --tags)
BRANCH   = $(shell git rev-parse --abbrev-ref HEAD)

STAGING_GCP_PROJECT_ID    := <set your staging gcp project id>
PRODUCTION_GCP_PROJECT_ID := <set your production gcp project id>

init: ## Install development tools
	@mkdir -p bin 2>/dev/null
	@go build -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
	@go mod download
.PHONY: init

lint: ## Run linters by golanci-lint
	@bin/golangci-lint run ./...
.PHONY: lint

test: ## Run simple test
	@go test -v -cover -race ./...
.PHONY: test

serve: ## Serve API service with .env file
	@go run .
.PHONY: server

deploy-stg: ## Deploy to STAGING environment
	@gcloud app deploy -q --no-promote --project=$(STAGING_GCP_PROJECT_ID) --version=$(BRANCH) staging.yaml
.PHONY: deploy-stg

deploy-prod: ## Deploy to PRODUCTION environment
	@gcloud app deploy -q --no-promote --project=$(PRODUCTION_GCP_PROJECT_ID) --version=$(TAG) production.yaml
.PHONY: deploy-prod

help: ## Show this help
	@perl -nle 'BEGIN {printf "Usage:\n  make \033[33m<target>\033[0m\n\nTargets:\n"} \
	printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 if /^([a-zA-Z_-].+):\s+## (.*)/' $(MAKEFILE_LIST)
.PHONY: help

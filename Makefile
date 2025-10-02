.PHONY: help test coverage bench lint fmt clean build install release version

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test ./... -v -race

coverage: ## Generate coverage report
	@echo "📊 Generating coverage..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage saved to coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test ./... -bench=. -benchmem -run=^$$

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else go vet ./...; fi

fmt: ## Format code
	@echo "✨ Formatting code..."
	@go fmt ./...

build: ## Build examples
	@echo "🔨 Building..."
	@mkdir -p bin
	@go build -o bin/basic ./examples/basic 2>/dev/null || echo "⚠️  Skipped"
	@go build -o bin/capacity-demo ./examples/capacity-demo 2>/dev/null || echo "⚠️  Skipped"

clean: ## Clean files
	@rm -rf bin/ coverage.out coverage.html
	@go clean

install: ## Install deps
	@go mod download
	@go mod tidy

version: ## Show version
	@git describe --tags --always 2>/dev/null || echo "v0.0.0"

release-patch: ## Patch release
	@$(MAKE) release BUMP=patch

release-minor: ## Minor release
	@$(MAKE) release BUMP=minor

release-major: ## Major release
	@$(MAKE) release BUMP=major

release: lint test ## Create release
	@CURRENT=$$(git describe --tags --always 2>/dev/null || echo "v0.0.0"); \
	if [ "$(BUMP)" = "major" ]; then \
	NEW=$$(echo $$CURRENT | sed 's/v//' | awk -F. '{print "v" $$1+1 ".0.0"}'); \
	elif [ "$(BUMP)" = "minor" ]; then \
	NEW=$$(echo $$CURRENT | sed 's/v//' | awk -F. '{print "v" $$1 "." $$2+1 ".0"}'); \
	elif [ "$(BUMP)" = "patch" ]; then \
	NEW=$$(echo $$CURRENT | sed 's/v//' | awk -F. '{print "v" $$1 "." $$2 "." $$3+1}'); \
	else \
	NEW="v1.0.0"; \
	fi; \
	echo "Version: $$CURRENT -> $$NEW"; \
	if [ -f CHANGELOG.md ]; then \
	sed -i "s/## \[Unreleased\]/## [Unreleased]\n\n## [$$NEW] - $$(date +%Y-%m-%d)/" CHANGELOG.md; \
	fi; \
	git add .; \
	git commit -m "chore: release $$NEW" || true; \
	git tag -a $$NEW -m "Release $$NEW"; \
	git push origin main; \
	git push origin $$NEW; \
	echo "✓ Released $$NEW"

ci: lint test bench ## CI checks
	@echo "✓ CI passed"

.DEFAULT_GOAL := help

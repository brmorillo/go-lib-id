.PHONY: help test coverage bench lint fmt clean build install release version next-version check-commits tag-version semantic-release setup-dev

help: ## Show help
	@echo "🆔 go-lib-id - Makefile Commands"
	@echo "================================="
	@echo ""
	@echo "📦 Development:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(test|coverage|bench|lint|fmt|build|clean|install)"
	@echo ""
	@echo "🚀 Versioning & Release:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(version|release|tag|semantic|setup|check)"
	@echo ""
	@echo "🔧 CI/CD:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(ci)"
	@echo ""
	@echo "💡 Quick Start:"
	@echo "  make setup-dev        # Setup versioning tools"
	@echo "  make test             # Run tests"
	@echo "  make version          # Check current version"
	@echo "  make next-version     # Preview next version"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  docs/VERSIONING.md           # Versioning guide"
	@echo "  docs/CONVENTIONAL_COMMITS.md # Commit format"
	@echo "  scripts/setup-versioning.sh  # Setup script"

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

version: ## Show current version
	@git describe --tags --always 2>/dev/null || echo "v0.0.0"

next-version: ## Show what the next version would be
	@echo "Analyzing commits for next version..."
	@if command -v npx >/dev/null 2>&1; then \
		npx semantic-release --dry-run --no-ci 2>/dev/null | grep -E "next release version|No release" || echo "Run: npm install -g semantic-release"; \
	else \
		echo "Install semantic-release: npm install -g semantic-release"; \
	fi

check-commits: ## Validate commit messages format
	@echo "Checking commit message format..."
	@if command -v npx >/dev/null 2>&1; then \
		npx commitlint --from=HEAD~10 --to=HEAD || echo "Install commitlint: npm install -g @commitlint/cli @commitlint/config-conventional"; \
	else \
		echo "Install commitlint: npm install -g @commitlint/cli @commitlint/config-conventional"; \
	fi

# Manual versioning (fallback when automatic fails)
tag-version: ## Create version tag manually (usage: make tag-version VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag-version VERSION=v1.0.0"; exit 1; fi
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✓ Tag $(VERSION) created and pushed"

release-patch: ## Patch release (automated)
	@$(MAKE) release BUMP=patch

release-minor: ## Minor release (automated)
	@$(MAKE) release BUMP=minor

release-major: ## Major release (automated)
	@$(MAKE) release BUMP=major

# Semantic release (requires semantic-release installed)
semantic-release: ## Run semantic-release locally
	@echo "Running semantic-release..."
	@if command -v npx >/dev/null 2>&1; then \
		npx semantic-release; \
	else \
		echo "Install semantic-release: npm install -g semantic-release"; \
	fi

# Setup development environment
setup-dev: ## Setup development environment with versioning tools
	@echo "Setting up development environment..."
	@if command -v npm >/dev/null 2>&1; then \
		echo "Installing semantic-release and commitlint..."; \
		npm install -g semantic-release@21 \
			@semantic-release/changelog@6 \
			@semantic-release/git@10 \
			@semantic-release/github@9 \
			@semantic-release/exec@6 \
			conventional-changelog-conventionalcommits@6 \
			@commitlint/cli@17 \
			@commitlint/config-conventional@17; \
		echo "✓ Development tools installed"; \
	else \
		echo "Node.js/npm not found. Please install Node.js first."; \
	fi

release: lint test ## Create release (automated versioning)
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

ci: lint test bench check-commits ## CI checks with commit validation
	@echo "✓ CI passed"

.DEFAULT_GOAL := help

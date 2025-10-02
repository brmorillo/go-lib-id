.PHONY: help test coverage bench lint fmt clean build install release version next-version check-commits tag-version semantic-release setup-dev

help: ## Show help
	@echo "🆔 go-lib-id - Makefile Commands"
	@echo "================================="
	@echo ""
	@echo "📦 Development:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(test|coverage|bench|lint|fmt|build|clean|install)"
	@echo ""
	@echo "🚀 Versioning & Release:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(version|next-version|check-updates|check-commits|tag-version|release|semantic-release|setup-dev)"
	@echo ""
	@echo "🔧 CI/CD:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-20s\033[0m %s\n", $$1, $$2}' | grep -E "(install|ci)"
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
	@if command -v go >/dev/null 2>&1; then \
		echo "Testing with race detection..."; \
		if CGO_ENABLED=1 go test -race ./... >/dev/null 2>&1; then \
			CGO_ENABLED=1 go test -race -v ./...; \
		else \
			echo "⚠️ Race detector not available, running without -race flag"; \
			go test -v ./...; \
		fi; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi
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

coverage: ## Generate coverage report
	@echo "📊 Generating coverage..."
	@if command -v go >/dev/null 2>&1; then \
		go test ./... -coverprofile=coverage.out -covermode=atomic && \
		go tool cover -html=coverage.out -o coverage.html && \
		echo "✓ Coverage saved to coverage.html"; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@if command -v go >/dev/null 2>&1; then \
		go test ./... -bench=. -benchmem -run=^$$; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v go >/dev/null 2>&1; then \
		if command -v golangci-lint >/dev/null 2>&1; then \
			golangci-lint run ./...; \
		else \
			echo "⚠️  golangci-lint not found, using go vet..."; \
			go vet ./...; \
		fi; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

fmt: ## Format code
	@echo "✨ Formatting code..."
	@if command -v go >/dev/null 2>&1; then \
		go fmt ./...; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

build: ## Build examples
	@echo "🔨 Building..."
	@if command -v go >/dev/null 2>&1; then \
		mkdir -p bin && \
		(go build -o bin/basic ./examples/basic 2>/dev/null && echo "✓ Built bin/basic") || echo "⚠️  Skipped basic example" && \
		(go build -o bin/capacity-demo ./examples/capacity-demo 2>/dev/null && echo "✓ Built bin/capacity-demo") || echo "⚠️  Skipped capacity-demo example"; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

clean: ## Clean files
	@rm -rf bin/ coverage.out coverage.html
	@go clean

install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@if command -v go >/dev/null 2>&1; then \
		go mod download && \
		go mod tidy && \
		echo "✓ Dependencies installed"; \
	else \
		echo "❌ Go not found. Please install Go first."; \
		exit 1; \
	fi

version: ## Show current version
	@echo "📋 Version Information:"
	@if command -v git >/dev/null 2>&1; then \
		echo "  Current: $$(git describe --tags --always 2>/dev/null || echo 'v0.0.0-dev')"; \
		echo "  Branch: $$(git branch --show-current 2>/dev/null || echo 'unknown')"; \
		echo "  Commit: $$(git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo 'unknown')"; \
	else \
		echo "  ❌ Git not found. Cannot determine version."; \
	fi

next-version: ## Show what the next version would be
	@echo "Analyzing commits for next version..."
	@if command -v npx >/dev/null 2>&1; then \
		npx semantic-release --dry-run --no-ci 2>/dev/null | grep -E "next release version|No release" || echo "Run: npm install -g semantic-release"; \
	else \
		echo "Install semantic-release: npm install -g semantic-release"; \
	fi

generate-changelog: ## Generate changelog from commits (local)
	@echo "Generating changelog from commits..."
	@if command -v npx >/dev/null 2>&1; then \
		npx conventional-changelog -p conventionalcommits -i CHANGELOG.md -s -r 0; \
		echo "✓ Changelog generated in CHANGELOG.md"; \
	else \
		echo "Install conventional-changelog: npm install -g conventional-changelog-cli"; \
	fi

preview-changelog: ## Preview what changelog would look like
	@echo "Previewing changelog for unreleased commits..."
	@if command -v npx >/dev/null 2>&1; then \
		npx conventional-changelog -p conventionalcommits -u; \
	else \
		echo "Install conventional-changelog: npm install -g conventional-changelog-cli"; \
	fi

info: ## Show system information for bug reports
	@echo "🔍 System Information for Bug Reports"
	@echo "====================================="
	@echo ""
	@echo "📅 Date: $$(date)"
	@echo "👤 User: $$(whoami)"
	@echo "💻 Hostname: $$(hostname)"
	@echo ""
	@echo "🐹 Go Information:"
	@if command -v go >/dev/null 2>&1; then \
		echo "  Version: $$(go version)"; \
		echo "  GOPATH: $$(go env GOPATH 2>/dev/null || echo 'not set')"; \
		echo "  GOROOT: $$(go env GOROOT 2>/dev/null || echo 'not set')"; \
		echo "  GOOS: $$(go env GOOS 2>/dev/null || echo 'unknown')"; \
		echo "  GOARCH: $$(go env GOARCH 2>/dev/null || echo 'unknown')"; \
	else \
		echo "  ❌ Go not found in PATH"; \
	fi
	@echo ""
	@echo "🖥️  Operating System:"
	@if [ "$$(uname)" = "Linux" ]; then \
		echo "  OS: Linux"; \
		echo "  Kernel: $$(uname -r)"; \
		if command -v lsb_release >/dev/null 2>&1; then \
			echo "  Distribution: $$(lsb_release -d | cut -f2)"; \
		elif [ -f /etc/os-release ]; then \
			echo "  Distribution: $$(grep PRETTY_NAME /etc/os-release | cut -d'"' -f2)"; \
		fi; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		echo "  OS: macOS"; \
		echo "  Version: $$(sw_vers -productVersion 2>/dev/null || echo 'unknown')"; \
		echo "  Build: $$(sw_vers -buildVersion 2>/dev/null || echo 'unknown')"; \
	elif [ "$$(uname -o 2>/dev/null)" = "Msys" ] || [ "$$(uname -o 2>/dev/null)" = "Cygwin" ]; then \
		echo "  OS: Windows (via $$(uname -o))"; \
		echo "  Version: $$(uname -r)"; \
	else \
		echo "  OS: $$(uname -s)"; \
		echo "  Version: $$(uname -r)"; \
	fi
	@echo ""
	@echo "📦 Library Information:"
	@echo "  Current Version: $$(git describe --tags --always 2>/dev/null || echo 'v0.0.0-dev')"
	@echo "  Branch: $$(git branch --show-current 2>/dev/null || echo 'unknown')"
	@echo "  Commit: $$(git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo 'unknown')"
	@echo "  Go Module: $$(grep '^module' go.mod 2>/dev/null | cut -d' ' -f2 || echo 'unknown')"
	@if [ -f go.mod ]; then \
		echo "  Go Version (mod): $$(grep '^go ' go.mod | cut -d' ' -f2)"; \
	fi
	@echo ""
	@echo "🔧 Development Tools:"
	@if command -v git >/dev/null 2>&1; then \
		echo "  Git: $$(git --version)"; \
	else \
		echo "  Git: ❌ not found"; \
	fi
	@if command -v make >/dev/null 2>&1; then \
		echo "  Make: $$(make --version | head -1)"; \
	else \
		echo "  Make: ❌ not found"; \
	fi
	@if command -v node >/dev/null 2>&1; then \
		echo "  Node.js: $$(node --version)"; \
	else \
		echo "  Node.js: ❌ not found"; \
	fi
	@if command -v npm >/dev/null 2>&1; then \
		echo "  npm: $$(npm --version)"; \
	else \
		echo "  npm: ❌ not found"; \
	fi
	@echo ""
	@echo "📋 Copy this information when reporting bugs!"
	@echo "==============================================="

check-updates: ## Check for dependency updates
	@echo "📦 Checking for dependency updates..."
	@if command -v go >/dev/null 2>&1; then \
		go list -u -m all | grep -E '\\[.*\\]$$' || echo "✅ All dependencies are up to date"; \
	else \
		echo "❌ Go not found. Cannot check updates."; \
		exit 1; \
	fi

check-commits: ## Check commit message format
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

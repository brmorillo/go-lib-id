
#!/bin/bash

# Setup script for automated versioning system
# This script configures all necessary tools for semantic versioning

set -e

echo "🚀 Setting up automated versioning system..."
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Node.js is installed
check_nodejs() {
    print_status "Checking Node.js installation..."
    if command -v node >/dev/null 2>&1; then
        NODE_VERSION=$(node --version)
        print_success "Node.js found: $NODE_VERSION"
        return 0
    else
        print_error "Node.js not found. Please install Node.js first:"
        echo "  - Visit: https://nodejs.org/"
        echo "  - Or use your package manager:"
        echo "    - Ubuntu/Debian: sudo apt install nodejs npm"
        echo "    - macOS: brew install node"
        echo "    - Windows: chocolatey install nodejs"
        return 1
    fi
}

# Install semantic-release and related tools
install_semantic_tools() {
    print_status "Installing semantic-release and commitlint..."
    
    npm install -g \
        semantic-release@21 \
        @semantic-release/changelog@6 \
        @semantic-release/git@10 \
        @semantic-release/github@9 \
        @semantic-release/exec@6 \
        conventional-changelog-conventionalcommits@6 \
        @commitlint/cli@17 \
        @commitlint/config-conventional@17
    
    print_success "Semantic tools installed successfully!"
}

# Setup git hooks for commitlint
setup_git_hooks() {
    print_status "Setting up git hooks for commit message validation..."
    
    # Create .husky directory if it doesn't exist
    mkdir -p .husky
    
    # Create commit-msg hook
    cat > .husky/commit-msg << 'EOF'
#!/bin/sh
. "$(dirname "$0")/_/husky.sh"

npx --no -- commitlint --edit ${1}
EOF
    
    chmod +x .husky/commit-msg
    
    # Install husky if not present
    if ! command -v npx husky >/dev/null 2>&1; then
        print_status "Installing husky for git hooks..."
        npm install -g husky@8
    fi
    
    print_success "Git hooks configured!"
}

# Validate current repository
validate_repo() {
    print_status "Validating repository setup..."
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        print_error "Not in a git repository!"
        return 1
    fi
    
    # Check if remote origin exists
    if ! git remote get-url origin >/dev/null 2>&1; then
        print_warning "No remote 'origin' found. Make sure to add your GitHub repository:"
        echo "  git remote add origin https://github.com/username/repo.git"
    fi
    
    print_success "Repository validation passed!"
}

# Test the setup
test_setup() {
    print_status "Testing the setup..."
    
    # Test semantic-release
    if command -v npx semantic-release >/dev/null 2>&1; then
        print_success "semantic-release is working"
    else
        print_error "semantic-release not found"
        return 1
    fi
    
    # Test commitlint
    if command -v npx commitlint >/dev/null 2>&1; then
        print_success "commitlint is working"
    else
        print_error "commitlint not found"
        return 1
    fi
    
    # Test our Makefile targets
    if make --dry-run version >/dev/null 2>&1; then
        print_success "Makefile targets are working"
    else
        print_warning "Some Makefile targets may not work"
    fi
}

# Show next steps
show_next_steps() {
    echo ""
    echo "🎉 Setup completed successfully!"
    echo "================================"
    echo ""
    echo "Next steps:"
    echo ""
    echo "1. 📝 Commit your changes using conventional format:"
    echo "   git add ."
    echo "   git commit -m \"feat: setup automated versioning system\""
    echo ""
    echo "2. 🚀 Push to trigger first automated release:"
    echo "   git push origin main"
    echo ""
    echo "3. 🔍 Available commands:"
    echo "   make version          # Show current version"
    echo "   make next-version     # Preview next version"
    echo "   make check-commits    # Validate commit messages"
    echo "   make semantic-release # Run release locally"
    echo ""
    echo "4. 📚 Documentation:"
    echo "   - Conventional Commits: docs/CONVENTIONAL_COMMITS.md"
    echo "   - Versioning Guide: docs/VERSIONING.md"
    echo ""
    echo "5. 🔗 GitHub Setup:"
    echo "   - Ensure GitHub Actions are enabled"
    echo "   - Check repository permissions for releases"
    echo "   - Consider enabling branch protection on main"
    echo ""
    print_success "Happy releasing! 🚀"
}

# Main execution
main() {
    validate_repo || exit 1
    check_nodejs || exit 1
    install_semantic_tools || exit 1
    setup_git_hooks || print_warning "Git hooks setup failed, but continuing..."
    test_setup || print_warning "Some tests failed, but setup might still work"
    show_next_steps
}

# Run if executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
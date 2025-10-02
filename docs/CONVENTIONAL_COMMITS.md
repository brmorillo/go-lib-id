# Conventional Commits Guide

This project follows the [Conventional Commits](https://www.conventionalcommits.org/) specification for automated versioning.

## Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Types

- **feat**: A new feature (triggers MINOR version bump)
- **fix**: A bug fix (triggers PATCH version bump)
- **docs**: Documentation only changes (triggers PATCH version bump)
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- **refactor**: A code change that neither fixes a bug nor adds a feature (triggers PATCH version bump)
- **perf**: A code change that improves performance (triggers PATCH version bump)
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files
- **revert**: Reverts a previous commit (triggers PATCH version bump)

## Breaking Changes

Add `BREAKING CHANGE:` in the footer or use `!` after the type/scope to trigger a MAJOR version bump.

## Examples

### Feature (Minor Release)
```
feat: add UUID v7 support

Implements RFC 9562 compliant UUID v7 generation with millisecond precision
timestamp and monotonic sequence for better database performance.
```

### Bug Fix (Patch Release)
```
fix: resolve race condition in sequence generation

The sequence counter was not properly protected in concurrent environments,
causing potential ID collisions under high load.

Fixes #123
```

### Breaking Change (Major Release)
```
feat!: change New() function signature

BREAKING CHANGE: The New() function now requires both processID and workerID parameters.
Previously, workerID was optional and defaulted to 0.

Migration guide:
- Old: idgen.New(5)
- New: idgen.New(5, 0)
```

### Documentation (Patch Release)
```
docs: add performance benchmarks to README

Include benchmark results and capacity analysis for different
configuration scenarios.
```

### Chore (No Release)
```
chore: update CI to use Go 1.21

Updates GitHub Actions workflow to use the latest Go version
for better performance and security.
```

## Scopes (Optional)

You can add scopes to provide more context:

- `api`: Changes to public API
- `core`: Changes to core functionality  
- `docs`: Documentation changes
- `test`: Test-related changes
- `ci`: CI/CD related changes

Examples:
```
feat(api): add batch generation methods
fix(core): resolve timestamp overflow
docs(api): improve method documentation
test(core): add concurrent generation tests
```

## Automated Versioning

Based on commit messages, versions are automatically determined:

- **MAJOR** (1.0.0 → 2.0.0): Breaking changes (`feat!:` or `BREAKING CHANGE:`)
- **MINOR** (1.0.0 → 1.1.0): New features (`feat:`)
- **PATCH** (1.0.0 → 1.0.1): Bug fixes, docs, refactors (`fix:`, `docs:`, `refactor:`, etc.)

## GitHub Releases

Releases are automatically created with:
- Changelog generated from commit messages
- Semantic version tags
- Release notes categorized by change type
- Assets (if applicable)

## Tools

- **semantic-release**: Automates the release process
- **commitlint**: Validates commit message format
- **conventional-changelog**: Generates changelogs

## Tips

1. Use present tense: "add feature" not "added feature"
2. Don't capitalize first letter: "add feature" not "Add feature"  
3. No period at the end: "add feature" not "add feature."
4. Be descriptive but concise
5. Reference issues when applicable: "Fixes #123" or "Closes #456"
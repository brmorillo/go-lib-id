# Automated Versioning System

This project uses **semantic versioning** with **automated releases** based on **conventional commits**.

## 🚀 How it Works

### 1. Commit Messages Drive Versioning

Every commit message determines the next version:

```bash
# Patch version (1.0.0 → 1.0.1)
git commit -m "fix: resolve race condition in ID generation"
git commit -m "docs: update README examples"
git commit -m "refactor: improve error handling"

# Minor version (1.0.0 → 1.1.0)  
git commit -m "feat: add UUID v7 support"
git commit -m "feat: implement batch generation"

# Major version (1.0.0 → 2.0.0)
git commit -m "feat!: change API signature"
git commit -m "feat: remove deprecated methods

BREAKING CHANGE: Removed NewSnowflake() function, use New() instead"
```

### 2. Automatic Process

When you push to `main`/`master`:

1. **GitHub Actions** triggers
2. **Tests run** and must pass
3. **Commits analyzed** for version bump
4. **Version calculated** using semantic rules
5. **Changelog updated** automatically
6. **Git tag created** (e.g., `v1.2.0`)
7. **GitHub release published** with notes
8. **Assets uploaded** (if configured)

### 3. Development Releases

Pushes to `dev` branch create **pre-releases**:
- Format: `v0.0.0-dev.123.abc1234`
- Marked as pre-release on GitHub
- Good for testing before main release

## 📋 Commit Message Rules

### Required Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types and Version Impact

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat` | New feature | **MINOR** | `feat: add ULID support` |
| `fix` | Bug fix | **PATCH** | `fix: resolve memory leak` |
| `docs` | Documentation | **PATCH** | `docs: improve API examples` |
| `refactor` | Code refactoring | **PATCH** | `refactor: simplify generator logic` |
| `perf` | Performance improvement | **PATCH** | `perf: optimize ID generation` |
| `test` | Tests | **NONE** | `test: add benchmark tests` |
| `chore` | Maintenance | **NONE** | `chore: update dependencies` |
| `ci` | CI/CD changes | **NONE** | `ci: add coverage reporting` |
| `build` | Build system | **NONE** | `build: update Go version` |

### Breaking Changes (MAJOR bump)

Add `!` after type or `BREAKING CHANGE:` in footer:

```bash
# Method 1: Exclamation mark
git commit -m "feat!: change New() signature"

# Method 2: Breaking change footer
git commit -m "feat: improve API

BREAKING CHANGE: The New() function now requires explicit workerID parameter"
```

## 🛠️ Manual Commands

### Check Current Version
```bash
make version
```

### Preview Next Version
```bash
make next-version
```

### Validate Commits
```bash
make check-commits
```

### Manual Tag (Emergency)
```bash
make tag-version VERSION=v1.2.3
```

### Setup Dev Environment
```bash
make setup-dev
```

## 🔧 Configuration Files

### `.releaserc.json`
Configures semantic-release behavior:
- Which branches trigger releases
- How commits map to version bumps
- Changelog generation rules
- GitHub release settings

### `.commitlintrc.json`
Validates commit message format:
- Enforces conventional commits
- Sets message length limits
- Defines allowed types and scopes

### `.github/workflows/release.yml`
GitHub Actions workflow:
- Runs tests before release
- Executes semantic-release
- Creates GitHub releases
- Handles dev branch pre-releases

## 📈 Release Process

### For Contributors

1. **Follow commit convention**:
   ```bash
   git commit -m "feat: add new ID type support"
   ```

2. **Push to dev** for testing:
   ```bash
   git push origin dev
   ```

3. **Create PR to main** when ready

### For Maintainers

1. **Merge PR** to main/master
2. **GitHub Actions** handles the rest automatically
3. **Release appears** on GitHub with changelog

## 🎯 Examples

### Adding a Feature
```bash
# Development
git checkout dev
git commit -m "feat: implement KSUID generator"
git push origin dev
# → Creates v0.0.0-dev.45.abc1234

# Production  
git checkout main
git merge dev
git push origin main
# → Creates v1.1.0 (minor bump)
```

### Bug Fix
```bash
git commit -m "fix: resolve sequence overflow issue

The sequence counter was overflowing after 4096 IDs per millisecond,
causing duplicate IDs in high-throughput scenarios.

Fixes #123"
git push origin main
# → Creates v1.0.1 (patch bump)
```

### Breaking Change
```bash
git commit -m "feat!: simplify API interface

BREAKING CHANGE: Removed deprecated NewSnowflake() function.
Use New() instead.

Migration:
- Old: generator, err := idgen.NewSnowflake(5, 12)  
- New: generator, err := idgen.New(5, 12)"
git push origin main
# → Creates v2.0.0 (major bump)
```

## 🚨 Troubleshooting

### "No release" Message
- Check commit messages follow convention
- Ensure commits have releasable types (`feat`, `fix`, etc.)
- Verify you're on the correct branch

### Release Failed
- Check GitHub Actions logs
- Verify GitHub token permissions
- Ensure tests are passing

### Wrong Version Bump
- Review commit message types
- Use `make next-version` to preview
- Consider if breaking change markers are needed

## 📚 References

- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [semantic-release](https://github.com/semantic-release/semantic-release)
- [Keep a Changelog](https://keepachangelog.com/)

## 🆘 Emergency Manual Release

If automation fails, you can create a manual release:

```bash
# 1. Create and push tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# 2. Create GitHub release manually
# Go to GitHub → Releases → Create new release → Select tag

# 3. Update CHANGELOG.md manually if needed
```

Remember: The automated system is designed to prevent human error and ensure consistency. Only use manual releases in emergencies!
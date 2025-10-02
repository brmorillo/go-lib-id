# Contributing to go-lib-id

Thank you for your interest in contributing to **go-lib-id**! 🎉

This is a public Go library for generating distributed unique IDs. We appreciate contributions, suggestions, bug reports, and improvements from the community.

## 🚀 How to Contribute

### 1. Report Issues

Found a bug or have an idea for improvement?

- Open an [issue](https://github.com/brmorillo/go-lib-id/issues) describing the problem or suggestion
- Be clear and specific in your description
- Include code examples when possible

### 2. Suggest New Features

Want to see a new ID type implemented?

- Check if it's already in the [roadmap](README.md#-available-id-comparison)
- Open an issue with the tag `enhancement`
- Describe the use case and benefits

### 3. Submit Pull Requests

Ready to contribute code?

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/YOUR-USER/go-lib-id.git`
3. Create a **branch** for your feature: `git checkout -b feature/my-feature`
4. **Implement** your changes with tests
5. **Run tests**: `go test ./...`
6. **Commit** with clear messages: `git commit -m "feat: add ULID implementation"`
7. **Push**: `git push origin feature/my-feature`
8. Open a **Pull Request** on GitHub

## 📋 Development Guidelines

### Code Standards

- Follow Go best practices and conventions
- Use `gofmt` to format code
- Add comprehensive tests for all new features
- Maintain 100% test coverage for critical code
- Include benchmarks for performance-critical code

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) standard:

- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation changes
- `test:` Adding or updating tests
- `refactor:` Code refactoring
- `perf:` Performance improvements

Examples:
```bash
feat: implement ULID generator
fix: race condition in UUID v7 generation
docs: update README with new examples
test: add concurrency tests for Snowflake
```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test ./... -bench=. -benchmem

# Run with race detection
go test ./... -race
```

### Documentation

- Document all exported functions and types
- Include usage examples in documentation comments
- Update README.md when adding new features
- Keep CHANGELOG.md updated

## 🎯 Priority Areas

We're actively seeking contributions in:

1. **New ID Types**: ULID, KSUID, xid, CUID, NanoID, ShortID, Sonyflake
2. **Performance Optimizations**: Benchmarks and improvements
3. **Documentation**: Examples and best practices guides
4. **Tests**: Edge cases and concurrency scenarios
5. **Bug Fixes**: Any issues found in production

## 📝 Publishing New Versions

### For Maintainers

When publishing a new version:

1. **Update CHANGELOG.md** with changes
2. **Run all tests**: `go test ./... -v`
3. **Run benchmarks**: `go test ./... -bench=.`
4. **Update version** in documentation if needed
5. **Commit changes**: `git commit -m "chore: prepare v1.x.x release"`
6. **Create tag**: `git tag v1.x.x`
7. **Push**: `git push origin main --tags`

Or use the automated script:

```bash
./publish.sh
```

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Breaking changes
- **MINOR** (v1.1.0): New features, backward compatible
- **PATCH** (v1.0.1): Bug fixes, backward compatible

## 🔍 Code Review Process

All Pull Requests will be reviewed by maintainers:

1. Code quality and standards
2. Test coverage
3. Documentation completeness
4. Performance impact
5. Security considerations

## 🌟 Recognition

Contributors will be recognized in:

- CHANGELOG.md (for significant contributions)
- GitHub contributors page
- Special mentions in releases

## 📞 Contact

Questions or suggestions?

- Open an [issue](https://github.com/brmorillo/go-lib-id/issues)
- Start a [discussion](https://github.com/brmorillo/go-lib-id/discussions)

## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for making **go-lib-id** better! 🙏


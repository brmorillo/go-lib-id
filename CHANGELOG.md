# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### 🚀 Features
- Refactor API from NewSnowflake() to New() for cleaner interface
- Add comprehensive documentation following Go conventions
- Move examples from cmd/* to examples/* following Go conventions
- Translate all examples and comments to English

### 📚 Documentation
- Add detailed examples and parameter descriptions for all methods
- Update README with new API usage and examples section
- Add professional-grade documentation for public library

### ♻️ Code Refactoring
- Maintain backward compatibility with deprecated NewSnowflake() functions
- Improve code organization and structure

### ✅ Tests
- All tests passing (35 tests) with new API
- Update test files to use new API naming

### 🏗️ Build System
- Update Makefile with new examples paths
- Add automated versioning with semantic-release
- Configure GitHub Actions for CI/CD

---

*Note: This changelog will be automatically maintained by semantic-release going forward.*

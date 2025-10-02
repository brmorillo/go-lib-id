# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.0.2](https://github.com/brmorillo/go-lib-id/compare/v1.0.1...v1.0.2) (2025-10-02)


### 🐛 Bug Fixes

* remove redundant make generate-changelog from prepareCmd ([b85f716](https://github.com/brmorillo/go-lib-id/commit/b85f716b0eca11ec8c0a802a7f0c6bd458c1f07b))

## [1.0.1](https://github.com/brmorillo/go-lib-id/compare/v1.0.0...v1.0.1) (2025-10-02)


### 🐛 Bug Fixes

* remove discussionCategoryName from release config (discussions not enabled) ([8121ba7](https://github.com/brmorillo/go-lib-id/commit/8121ba75d646455bf3a8b5a9f457427c3e2aa64f))

## 1.0.0 (2025-10-02)


### 🚀 Features

* enhance documentation and automated releases with professional multi-platform binaries ([8c48d41](https://github.com/brmorillo/go-lib-id/commit/8c48d41ca2f4d2a1272cf52276063221d4b4f98f))
* enhance Makefile with changelog generation and system info commands ([931c459](https://github.com/brmorillo/go-lib-id/commit/931c4590c6b7961b4bf1c09d225e26440490e736))
* implement automated versioning system with semantic-release ([a260617](https://github.com/brmorillo/go-lib-id/commit/a260617579fb319ff2e8da155f412b2441cac02b))
* refactor API and improve documentation ([cbf2b6f](https://github.com/brmorillo/go-lib-id/commit/cbf2b6ff48f7a3f499ce177e8fb10c8e97a3a80f))


### 🐛 Bug Fixes

* add missing @semantic-release/exec dependency to release workflow ([ada9421](https://github.com/brmorillo/go-lib-id/commit/ada942136b285e22b681501181c0dab4b36e2098))
* pipeline ([#3](https://github.com/brmorillo/go-lib-id/issues/3)) ([104a12d](https://github.com/brmorillo/go-lib-id/commit/104a12d1792074266f926e9b188b2e3f236b20ad))
* releaserc.json ([89afa84](https://github.com/brmorillo/go-lib-id/commit/89afa84e4395f5e94ca4d7d259883f3921fa73b0))
* resolve CI/CD compatibility issues and add local testing tools ([#1](https://github.com/brmorillo/go-lib-id/issues/1)) ([0023eb3](https://github.com/brmorillo/go-lib-id/commit/0023eb3731bd2773262eed7442d32520c8671862))
* simplify semantic-release configuration ([2976e49](https://github.com/brmorillo/go-lib-id/commit/2976e493af36f0039f1e5218673a307932140ff1))
* upload-artifact@v4 ([#2](https://github.com/brmorillo/go-lib-id/issues/2)) ([cad183e](https://github.com/brmorillo/go-lib-id/commit/cad183e91ca272183202ae292f5e2ffbad6d479f))


### ♻️ Code Refactoring

* reorganize CI/CD pipeline for proper dev/main workflow ([10c394b](https://github.com/brmorillo/go-lib-id/commit/10c394bf37eb512ddb5098aedb565d3ceb3ac226))

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

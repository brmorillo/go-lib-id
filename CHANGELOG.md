# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Nothing yet

### Changed
- Nothing yet

### Fixed
- Nothing yet

## [1.0.0] - 2025-10-02

### Added
- **Snowflake ID Generator** - Discord/Twitter format (64-bit)
  - Process ID (0-31) and Worker ID (0-31) support
  - ~4.1M IDs per millisecond capacity
  - Batch generation support
  - Component extraction (timestamp, process, worker, sequence)
  - Global API for simplified usage
  - Comprehensive test suite (100% coverage)
  - Performance: ~275 ns/op with 0 allocations

- **UUID v4 Generator** - Random UUIDs
  - RFC 4122 compliant
  - Cryptographically secure (crypto/rand)
  - Batch generation support
  - 122 bits of randomness

- **UUID v7 Generator** - Time-ordered UUIDs
  - RFC 9562 compliant (draft)
  - 48-bit timestamp + 12-bit sequence + 62-bit random
  - Monotonically increasing within same millisecond
  - Thread-safe with mutex
  - Better database index performance than UUID v4

- **Placeholder structures for future implementations:**
  - CUID (Collision-resistant Unique Identifier)
  - ULID (Universally Unique Lexicographically Sortable ID)
  - NanoID (Compact URL-safe ID)
  - ShortID (Compressed UUID)
  - KSUID (K-Sortable Unique Identifier)
  - xid (MongoDB-like ObjectID)
  - Sonyflake (Sony's Snowflake variant)

### Documentation
- Comprehensive README with comparison table
- Usage examples for all implemented IDs
- "When to use each ID type" guide
- Visual format comparison
- API documentation in code

### Testing
- 100% test coverage for implemented features
- Concurrency tests
- Uniqueness validation
- Time-ordering validation
- Benchmark tests included

[1.0.0]: https://github.com/brmorillo/go-lib-id/releases/tag/v1.0.0

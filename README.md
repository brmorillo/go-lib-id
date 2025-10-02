# go-lib-id

[![Go Referenc| **Sonyflake** | 63-bit | ✅ | Decimal | Improved Snowflake (174 years) | 🔴 |

## 🚀 Installation](https://pkg.go.dev/badge/github.com/brmorillo/go-lib-id.svg)](https://pkg.go.dev/github.com/brmorillo/go-lib-id)
[![Go Report Card](https://goreportcard.com/badge/github.com/brmorillo/go-lib-id)](https://goreportcard.com/report/github.com/brmorillo/go-lib-id)
[![CI](https://github.com/brmorillo/go-lib-id/actions/workflows/ci.yml/badge.svg)](https://github.com/brmorillo/go-lib-id/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/brmorillo/go-lib-id)](https://go.dev/)

🆔 Go library for generating unique and distributed IDs

## 📋 About

Complete library for generating different types of IDs for distributed systems, with support for multiple formats and unique identification strategies.

## 📊 ID Comparison

| ID Type       | Size      | Sortable | Encoding | Best For                       | Implemented |
| ------------- | --------- | -------- | -------- | ------------------------------ | ----------- |
| **Snowflake** | 64-bit    | ✅        | Decimal  | Numeric distributed IDs        | 🟢           |
| **UUID v4**   | 128-bit   | ❌        | Hex      | Maximum uniqueness             | 🟢           |
| **UUID v7**   | 128-bit   | ✅        | Hex      | Sortable UUID                  | 🟢           |
| **ULID**      | 128-bit   | ✅        | Base32   | URL-safe, case-insensitive     | 🔴           |
| **KSUID**     | 160-bit   | ✅        | Base62   | Distributed, second-precision  | 🔴           |
| **xid**       | 96-bit    | ✅        | Base32   | MongoDB-like                   | 🔴           |
| **CUID**      | ~25 chars | ✅        | Base36   | Collision-resistant            | 🔴           |
| **NanoID**    | 21 chars  | ❌        | Custom   | Short URLs                     | 🔴           |
| **ShortID**   | 22 chars  | ❌        | Base62   | Compact UUID                   | 🔴           |
| **Sonyflake** | 63-bit    | ✅        | Decimal  | Improved Snowflake (174 years) | �           |

## 🚀 Installation

```bash
go get github.com/brmorillo/go-lib-id
```

## 📖 Usage

### Snowflake ID

```go
package main

import (
    "fmt"
    "time"
    "github.com/brmorillo/go-lib-id/pkg/idgen"
)

func main() {
    // Create ID generator with process and worker IDs
    generator, err := idgen.New(5, 12)
    if err != nil {
        panic(err)
    }
    
    // Generate unique ID
    id := generator.Generate()
    fmt.Printf("Snowflake ID: %d\n", id)
    
    // Generate multiple IDs efficiently
    ids := generator.GenerateBatch(10)
    fmt.Printf("Generated %d IDs\n", len(ids))
    
    // Extract ID components for analysis
    timestamp := generator.ExtractTimestamp(id)
    processID := generator.ExtractProcessID(id)
    workerID := generator.ExtractWorkerID(id)
    sequence := generator.ExtractSequence(id)
    dateTime := generator.ExtractTime(id)
    
    fmt.Printf("Timestamp: %d, Process: %d, Worker: %d, Sequence: %d\n", 
        timestamp, processID, workerID, sequence)
    fmt.Printf("DateTime: %s\n", dateTime.Format(time.RFC3339))
}
```

### UUID v4

```go
package main

import (
    "fmt"
    "github.com/brmorillo/go-lib-id/pkg/idgen"
)

func main() {
    // Generate UUID v4
    uuid, err := idgen.GenerateUUIDv4()
    if err != nil {
        panic(err)
    }
    fmt.Printf("UUID v4: %s\n", uuid)
    
    // Generate multiple UUIDs
    uuids, err := idgen.GenerateUUIDv4Batch(10)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Generated %d UUIDs\n", len(uuids))
}
```

### UUID v7

```go
package main

import (
    "fmt"
    "github.com/brmorillo/go-lib-id/pkg/idgen"
)

func main() {
    // Generate UUID v7 (time-ordered)
    uuid, err := idgen.GenerateUUIDv7()
    if err != nil {
        panic(err)
    }
    fmt.Printf("UUID v7: %s\n", uuid)
    
    // Generate multiple UUID v7s
    uuids, err := idgen.GenerateUUIDv7Batch(10)
    if err != nil {
        panic(err)
    }
    
    // UUID v7s are ordered by creation time
    for i, id := range uuids {
        fmt.Printf("%d: %s\n", i+1, id)
    }
}
```

### Simplified Usage (Global API)

```go
package main

import (
    "fmt"
    "github.com/brmorillo/go-lib-id/pkg/idgen"
)

func main() {
    // Configure Process ID and Worker ID globally once
    // Recommendation: get from environment variables
    err := idgen.SetDefaultMachineID(10, 20)
    if err != nil {
        panic(err)
    }
    
    // Generate IDs directly
    id1 := idgen.GenerateSnowflake()
    id2 := idgen.GenerateSnowflake()
    
    fmt.Printf("ID 1: %d\n", id1)
    fmt.Printf("ID 2: %d\n", id2)
}
```

## 📁 Examples

The repository includes practical examples demonstrating library usage:

### Basic Usage Example
```bash
go run examples/basic/main.go
```
Comprehensive demo showing:
- Snowflake generator creation
- Individual and batch ID generation
- Component extraction from IDs
- Global API usage
- Distributed system simulation
- Performance testing
- Uniqueness verification

### Capacity Demo
```bash
go run examples/capacity-demo/main.go
```
Performance demonstration showing:
- Single millisecond capacity testing
- Continuous generation benchmarks
- Multi-worker simulation
- Theoretical vs actual performance
- Scalability analysis

Both examples include detailed output explaining each step and the theoretical limits of the Snowflake ID system.

## 🏗️ Snowflake Architecture

```
Snowflake ID (64 bits) - Discord/Twitter Format:
┌─────────┬──────────────┬────────────┬───────────┬──────────┐
│ Sign    │  Timestamp   │ Process ID │ Worker ID │ Sequence │
│ 1 bit   │   41 bits    │   5 bits   │  5 bits   │ 12 bits  │
│ (unused)│              │  (0-31)    │  (0-31)   │ (0-4095) │
└─────────┴──────────────┴────────────┴───────────┴──────────┘
```

## 🎯 When to Use Each ID Type

### Snowflake ❄️
**Use when:**
- Need numeric IDs (int64)
- Distributed system with multiple servers
- Temporal ordering is important
- Performance is critical (ultra-fast generation)
- Want IDs smaller than UUID

**Don't use when:**
- Cannot coordinate Process/Worker IDs
- Need more than 1024 simultaneous generators

### UUID v4 🎲
**Use when:**
- Need maximum randomness
- Creation order doesn't matter
- Compatibility with existing systems
- Don't want coordination between servers

**Don't use when:**
- Time ordering is important
- Need compact IDs
- Index performance is critical

### UUID v7 📅
**Use when:**
- Want UUID + time ordering
- Need UUID compatibility
- Database indexes are important
- Want better performance than UUID v4 on inserts

**Don't use when:**
- Need IDs more compact than 128 bits
- Want more friendly format than hex

### ULID 🔤 (Planned)
**Use when:**
- Want case-insensitive (URLs, emails)
- Need time ordering
- Want more readable format than UUID
- Base32 is preferable to hex

### KSUID ⏰ (Planned)
**Use when:**
- Second precision is sufficient
- Want more randomness than ULID
- Highly distributed system
- Want URL-safe format

### xid 🗄️ (Planned)
**Use when:**
- Want MongoDB-like compatibility
- IDs more compact than UUID (96 bits)
- Simple distributed system

### CUID 🛡️ (Planned)
**Use when:**
- Collision resistance is priority
- Offline/decentralized systems
- Want automatic machine identification

### NanoID �� (Planned)
**Use when:**
- Need short URLs
- 21 characters is sufficient
- Want to customize alphabet
- Public/user-facing IDs

### ShortID 📏 (Planned)
**Use when:**
- Want compressed UUID
- 22 characters is acceptable
- Base62 is suitable

### Sonyflake 🌸 (Planned)
**Use when:**
- Need more than 69 years of lifetime
- Want to support 65K+ machines
- 10ms precision is sufficient
- Improved Snowflake alternative

## 🎨 Visual Format Comparison

```
Snowflake:  1234567890123456789              (19 digits)
UUID v4:    550e8400-e29b-41d4-a716-446655440000  (36 chars)
UUID v7:    018b5e0c-3e4a-7000-8000-000000000000  (36 chars)
ULID:       01ARZ3NDEKTSV4RRFFQ69G5FAV       (26 chars)
KSUID:      0ujtsYcgvSTl8PAuAdqWYSMnLOv      (27 chars)
xid:        9m4e2mr0ui3e8a215n4g              (20 chars)
CUID:       cjld2cjxh0000qzrmn831i7rn         (~25 chars)
NanoID:     V1StGXR8_Z5jdHi6B-myT             (21 chars)
ShortID:    ppBqWA9fuP3FcvjJHQxNz3            (22 chars)
Sonyflake:  123456789012345                   (~15 digits)
```

## 🧪 Testing

Run all tests:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Run benchmarks:
```bash
go test ./... -bench=. -benchmem
```

## 📈 Performance

Snowflake ID generation benchmarks:
```
BenchmarkSnowflakeGenerate-12     4,362,130    275 ns/op    0 B/op    0 allocs/op
```

## 🤝 Contributing

Contributions are welcome! Please read our [contributing guidelines](CONTRIBUTING.md) and follow our [conventional commit format](docs/CONVENTIONAL_COMMITS.md) for automatic versioning.

### Quick Start for Contributors

1. **Fork and clone** the repository
2. **Create a feature branch**: `git checkout -b feat/amazing-feature`
3. **Make your changes** following the code style
4. **Write tests** for your changes
5. **Commit using conventional format**: `git commit -m "feat: add amazing feature"`
6. **Push to your fork**: `git push origin feat/amazing-feature`
7. **Open a Pull Request** with a clear description

### Commit Message Format

We use [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning:

```bash
# Feature (minor version bump)
git commit -m "feat: add ULID support"

# Bug fix (patch version bump)
git commit -m "fix: resolve race condition"

# Breaking change (major version bump)
git commit -m "feat!: change API signature"
```

See [docs/CONVENTIONAL_COMMITS.md](docs/CONVENTIONAL_COMMITS.md) for detailed guidelines.

## 🚀 Versioning

This project uses [Semantic Versioning](https://semver.org/) with automated releases:

- **Commits drive versioning**: Your commit messages determine version bumps
- **Automatic releases**: Push to main triggers GitHub Actions
- **Changelog generation**: Auto-generated from commit messages
- **GitHub releases**: Created automatically with release notes

See [docs/VERSIONING.md](docs/VERSIONING.md) for complete details.

### Current Version

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/brmorillo/go-lib-id)](https://github.com/brmorillo/go-lib-id/releases/latest)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/brmorillo/go-lib-id)](https://github.com/brmorillo/go-lib-id/tags)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Documentation](https://pkg.go.dev/github.com/brmorillo/go-lib-id)
- [Issues](https://github.com/brmorillo/go-lib-id/issues)
- [Releases](https://github.com/brmorillo/go-lib-id/releases)
- [Discussions](https://github.com/brmorillo/go-lib-id/discussions)

## 🌟 Acknowledgments

- Inspired by Twitter/Discord Snowflake
- UUID v7 based on RFC 9562 draft
- Thanks to the Go community


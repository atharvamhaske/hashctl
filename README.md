# hashctl ⟡

A beautiful terminal UI for computing cryptographic hashes.

![hashctl](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)
![Release](https://img.shields.io/github/v/release/atharvamhaske/hashctl?style=flat-square&include_prereleases=false)


## Features

- **Interactive TUI** — keyboard-driven interface with Bubble Tea
- **20+ algorithms** — SHA-256, SHA-512, BLAKE2, SHA-3, MD5, bcrypt, Argon2id...
- **Hash strings or files** — simple input modes
- **Clean aesthetic** — minimal, focused design

## Installation Guide

### Install via Go

```bash
go install github.com/atharvamhaske/hashctl@latest
```

### Build from Source

```bash
git clone https://github.com/atharvamhaske/hashctl
cd hashctl
go build -o hashctl .
```

### Commands

```bash
hashctl          # Launch TUI
hashctl list     # Show all algorithms
hashctl version  # Print version info
```

## Use as a Go pkg

You can import hashctl's hasher package in your own Go code:

```go
package main

import (
    "fmt"
    "github.com/atharvamhaske/hashctl/internal/hasher"
)

func main() {
    // Hash a string
    opts := hasher.DefaultOptions()
    opts.Algorithm = "sha256"
    
    result := hasher.HashString("hello world", opts)
    if result.Error != nil {
        panic(result.Error)
    }
    fmt.Println(result.Hash)
    // Output: b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9

    // Hash a file
    fileResult := hasher.HashFile("/path/to/file.txt", opts)
    if fileResult.Error != nil {
        panic(fileResult.Error)
    }
    fmt.Println(fileResult.Hash)

    // Hash multiple files in parallel
    files := []string{"file1.txt", "file2.txt", "file3.txt"}
    hasher.HashFiles(files, opts, func(r hasher.Result) {
        fmt.Printf("%s  %s\n", r.Hash, r.Input)
    })

    // List available algorithms
    for _, name := range hasher.ListNames() {
        alg, _ := hasher.GetAlgorithm(name)
        fmt.Printf("%s: %s\n", name, alg.Description)
    }
}
```

### List of all Functions

```go
// Hash a string
hasher.HashString(input string, opts Options) Result

// Hash a single file
hasher.HashFile(filename string, opts Options) Result

// Hash multiple files in parallel with ordered output
hasher.HashFiles(files []string, opts Options, onResult func(Result))

// Get algorithm by name
hasher.GetAlgorithm(name string) (Algorithm, bool)

// List all algorithm names
hasher.ListNames() []string

// Get algorithms grouped by category
hasher.GetAlgorithmsByCategory() map[Category][]Algorithm
```

## Available Algorithms

### Checksums
- CRC32

### Fast Cryptographic Hashes
- MD5, SHA-1
- SHA-224, SHA-256, SHA-384, SHA-512
- SHA-512/224, SHA-512/256
- SHA3-224, SHA3-256, SHA3-384, SHA3-512
- RIPEMD-160
- BLAKE2b-256, BLAKE2b-384, BLAKE2b-512
- BLAKE2s-256

### Password Hashing
- bcrypt
- Argon2id

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling
- [Cobra](https://github.com/spf13/cobra) — CLI

## License

[MIT](LICENSE)

# hashctl ⟡

A beautiful terminal UI for computing cryptographic hashes build using bubbletea, lipgloss in Go.

![hashctl](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)

## Features

- **Interactive TUI** — keyboard-driven interface with Bubble Tea
- **20+ algorithms** — SHA-256, SHA-512, BLAKE2, SHA-3, MD5, bcrypt, Argon2id...
- **Hash strings or files** — simple input modes
- **Clean aesthetic** — minimal, focused design

## Install

```bash
go install github.com/atharvamhaske/hashctl@latest
```

Or build from source:

```bash
git clone https://github.com/atharvamhaske/hashctl
cd hashctl
go build -o hashctl .
```

## Usage

Just run it:

```bash
hashctl
```

Use arrow keys to select an algorithm, then choose to hash a string or file.

### Commands

```bash
hashctl          # Launch TUI
hashctl list     # Show all algorithms
hashctl version  # Print version info
```

## Keyboard

| Key | Action |
|-----|--------|
| `↑` `k` | Move up |
| `↓` `j` | Move down |
| `enter` | Select / Confirm |
| `s` | Hash string |
| `f` | Hash file |
| `n` | New hash (same algorithm) |
| `r` | Restart |
| `esc` | Go back |
| `q` | Quit |

## Algorithms

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

MIT

// Package hasher provides cryptographic hash algorithms and utilities
package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/crc32"
	"sort"

	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Category represents the type of hashing algorithm
type Category int

const (
	CategoryChecksum Category = iota
	CategoryFastHash
	CategoryPasswordHash
)

func (c Category) String() string {
	switch c {
	case CategoryChecksum:
		return "Checksums (Non-Cryptographic)"
	case CategoryFastHash:
		return "Fast Cryptographic Hashes"
	case CategoryPasswordHash:
		return "Password Hashing / KDFs"
	default:
		return "Unknown"
	}
}

// Algorithm represents a hashing algorithm with metadata
type Algorithm struct {
	Name        string
	Description string
	Category    Category
	NewHash     func() hash.Hash
	// For password hashes, we use a different interface
	IsPasswordHash bool
}

// Registry holds all available algorithms
var Registry = map[string]Algorithm{
	// Checksums (Non-Cryptographic)
	"crc32": {
		Name:        "CRC32",
		Description: "Fast checksum for detecting accidental data corruption; not suitable for security.",
		Category:    CategoryChecksum,
		NewHash:     func() hash.Hash { return crc32.NewIEEE() },
	},

	// Fast Cryptographic Hashes
	"md5": {
		Name:        "MD5",
		Description: "128-bit hash, widely used but cryptographically broken. Use only for legacy compatibility.",
		Category:    CategoryFastHash,
		NewHash:     md5.New,
	},
	"sha1": {
		Name:        "SHA-1",
		Description: "160-bit hash, deprecated for security use. Common in legacy systems and git.",
		Category:    CategoryFastHash,
		NewHash:     sha1.New,
	},
	"sha224": {
		Name:        "SHA-224",
		Description: "Truncated variant of SHA-256 with 224-bit output.",
		Category:    CategoryFastHash,
		NewHash:     sha256.New224,
	},
	"sha256": {
		Name:        "SHA-256",
		Description: "Cryptographic hash widely used for integrity checks and content addressing.",
		Category:    CategoryFastHash,
		NewHash:     sha256.New,
	},
	"sha384": {
		Name:        "SHA-384",
		Description: "Truncated variant of SHA-512 with 384-bit output.",
		Category:    CategoryFastHash,
		NewHash:     sha512.New384,
	},
	"sha512": {
		Name:        "SHA-512",
		Description: "512-bit hash from the SHA-2 family, suitable for high-security applications.",
		Category:    CategoryFastHash,
		NewHash:     sha512.New,
	},
	"sha512-224": {
		Name:        "SHA-512/224",
		Description: "SHA-512 truncated to 224 bits, optimized for 64-bit platforms.",
		Category:    CategoryFastHash,
		NewHash:     sha512.New512_224,
	},
	"sha512-256": {
		Name:        "SHA-512/256",
		Description: "SHA-512 truncated to 256 bits, optimized for 64-bit platforms.",
		Category:    CategoryFastHash,
		NewHash:     sha512.New512_256,
	},
	"sha3-224": {
		Name:        "SHA3-224",
		Description: "224-bit SHA-3 hash based on Keccak sponge construction.",
		Category:    CategoryFastHash,
		NewHash:     sha3.New224,
	},
	"sha3-256": {
		Name:        "SHA3-256",
		Description: "256-bit SHA-3 hash, NIST standard alternative to SHA-256.",
		Category:    CategoryFastHash,
		NewHash:     sha3.New256,
	},
	"sha3-384": {
		Name:        "SHA3-384",
		Description: "384-bit SHA-3 hash based on Keccak sponge construction.",
		Category:    CategoryFastHash,
		NewHash:     sha3.New384,
	},
	"sha3-512": {
		Name:        "SHA3-512",
		Description: "512-bit SHA-3 hash, highest security level in SHA-3 family.",
		Category:    CategoryFastHash,
		NewHash:     sha3.New512,
	},
	"ripemd160": {
		Name:        "RIPEMD-160",
		Description: "160-bit hash used in Bitcoin addresses and PGP fingerprints.",
		Category:    CategoryFastHash,
		NewHash:     ripemd160.New,
	},
	"blake2b-256": {
		Name:        "BLAKE2b-256",
		Description: "Fast cryptographic hash, faster than MD5 while being secure.",
		Category:    CategoryFastHash,
		NewHash:     func() hash.Hash { h, _ := NewBlake2b256(); return h },
	},
	"blake2b-384": {
		Name:        "BLAKE2b-384",
		Description: "384-bit BLAKE2b variant, optimized for 64-bit platforms.",
		Category:    CategoryFastHash,
		NewHash:     func() hash.Hash { h, _ := NewBlake2b384(); return h },
	},
	"blake2b-512": {
		Name:        "BLAKE2b-512",
		Description: "512-bit BLAKE2b, one of the fastest secure hash functions.",
		Category:    CategoryFastHash,
		NewHash:     func() hash.Hash { h, _ := NewBlake2b512(); return h },
	},
	"blake2s-256": {
		Name:        "BLAKE2s-256",
		Description: "BLAKE2s optimized for 8-32 bit platforms and small inputs.",
		Category:    CategoryFastHash,
		NewHash:     func() hash.Hash { h, _ := NewBlake2s256(); return h },
	},

	// Password Hashing / KDFs
	"bcrypt": {
		Name:           "bcrypt",
		Description:    "Adaptive password hashing with configurable cost factor.",
		Category:       CategoryPasswordHash,
		IsPasswordHash: true,
	},
	"argon2id": {
		Name:           "Argon2id",
		Description:    "Memory-hard password hashing algorithm designed to resist GPU and ASIC attacks.",
		Category:       CategoryPasswordHash,
		IsPasswordHash: true,
	},
}

// GetAlgorithmsByCategory returns algorithms grouped by category
func GetAlgorithmsByCategory() map[Category][]Algorithm {
	result := make(map[Category][]Algorithm)
	for _, alg := range Registry {
		result[alg.Category] = append(result[alg.Category], alg)
	}
	// Sort each category by name
	for cat := range result {
		sort.Slice(result[cat], func(i, j int) bool {
			return result[cat][i].Name < result[cat][j].Name
		})
	}
	return result
}

// GetSortedAlgorithms returns all algorithms sorted by category then name
func GetSortedAlgorithms() []Algorithm {
	var algs []Algorithm
	for _, alg := range Registry {
		algs = append(algs, alg)
	}
	sort.Slice(algs, func(i, j int) bool {
		if algs[i].Category != algs[j].Category {
			return algs[i].Category < algs[j].Category
		}
		return algs[i].Name < algs[j].Name
	})
	return algs
}

// GetAlgorithm returns an algorithm by name (case-insensitive key)
func GetAlgorithm(name string) (Algorithm, bool) {
	alg, ok := Registry[name]
	return alg, ok
}

// ListNames returns all algorithm names
func ListNames() []string {
	var names []string
	for name := range Registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

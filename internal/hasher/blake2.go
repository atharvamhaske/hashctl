package hasher

import (
	"hash"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

// NewBlake2b256 creates a new BLAKE2b-256 hash
func NewBlake2b256() (hash.Hash, error) {
	return blake2b.New256(nil)
}

// NewBlake2b384 creates a new BLAKE2b-384 hash
func NewBlake2b384() (hash.Hash, error) {
	return blake2b.New384(nil)
}

// NewBlake2b512 creates a new BLAKE2b-512 hash
func NewBlake2b512() (hash.Hash, error) {
	return blake2b.New512(nil)
}

// NewBlake2s256 creates a new BLAKE2s-256 hash
func NewBlake2s256() (hash.Hash, error) {
	return blake2s.New256(nil)
}

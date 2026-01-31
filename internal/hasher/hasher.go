package hasher

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Result represents the result of a hash computation
type Result struct {
	Input    string // filename or string input
	Hash     string // hex-encoded hash
	Error    error  // any error that occurred
	IsFile   bool   // true if input is a file
	Duration time.Duration
}

// Options for hash computation
type Options struct {
	Algorithm   string
	Parallelism int
	// For password hashing
	BcryptCost   int
	Argon2Time   uint32
	Argon2Memory uint32
	Argon2Lanes  uint8
	Argon2KeyLen uint32
}

// DefaultOptions returns sensible defaults
func DefaultOptions() Options {
	return Options{
		Algorithm:    "sha256",
		Parallelism:  runtime.NumCPU(),
		BcryptCost:   bcrypt.DefaultCost,
		Argon2Time:   1,
		Argon2Memory: 64 * 1024, // 64MB
		Argon2Lanes:  4,
		Argon2KeyLen: 32,
	}
}

// HashString computes the hash of a string
func HashString(input string, opts Options) Result {
	start := time.Now()

	alg, ok := GetAlgorithm(opts.Algorithm)
	if !ok {
		return Result{
			Input: input,
			Error: fmt.Errorf("unknown algorithm: %s", opts.Algorithm),
		}
	}

	var hashStr string
	var err error

	if alg.IsPasswordHash {
		hashStr, err = hashPassword(input, opts)
	} else {
		h := alg.NewHash()
		h.Write([]byte(input))
		hashStr = hex.EncodeToString(h.Sum(nil))
	}

	return Result{
		Input:    input,
		Hash:     hashStr,
		Error:    err,
		IsFile:   false,
		Duration: time.Since(start),
	}
}

// HashFile computes the hash of a file using streaming
func HashFile(filename string, opts Options) Result {
	start := time.Now()

	alg, ok := GetAlgorithm(opts.Algorithm)
	if !ok {
		return Result{
			Input:  filename,
			Error:  fmt.Errorf("unknown algorithm: %s", opts.Algorithm),
			IsFile: true,
		}
	}

	// Password hashes require reading entire file into memory
	if alg.IsPasswordHash {
		data, err := os.ReadFile(filename)
		if err != nil {
			return Result{
				Input:  filename,
				Error:  err,
				IsFile: true,
			}
		}
		hashStr, err := hashPassword(string(data), opts)
		return Result{
			Input:    filename,
			Hash:     hashStr,
			Error:    err,
			IsFile:   true,
			Duration: time.Since(start),
		}
	}

	// Stream-based hashing for regular algorithms
	hashStr, err := computeFileHash(alg.NewHash(), filename)
	return Result{
		Input:    filename,
		Hash:     hashStr,
		Error:    err,
		IsFile:   true,
		Duration: time.Since(start),
	}
}

// HashFiles computes hashes for multiple files in parallel while preserving order
func HashFiles(files []string, opts Options, onResult func(Result)) {
	if len(files) == 0 {
		return
	}

	alg, ok := GetAlgorithm(opts.Algorithm)
	if !ok {
		for _, f := range files {
			onResult(Result{
				Input:  f,
				Error:  fmt.Errorf("unknown algorithm: %s", opts.Algorithm),
				IsFile: true,
			})
		}
		return
	}

	// Channel for limiting concurrent processing
	sem := make(chan struct{}, opts.Parallelism)

	// Index of the current file that will be printed
	currFile := int32(0)

	// Store results for ordered output
	results := make([]Result, len(files))
	printed := make([]bool, len(files))
	var printMu sync.Mutex

	wg := sync.WaitGroup{}
	wg.Add(len(files))

	for i, fname := range files {
		i := int32(i)
		fname := fname
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			start := time.Now()
			var hashStr string
			var err error

			if alg.IsPasswordHash {
				data, readErr := os.ReadFile(fname)
				if readErr != nil {
					err = readErr
				} else {
					hashStr, err = hashPassword(string(data), opts)
				}
			} else {
				h := alg.NewHash()
				hashStr, err = computeFileHash(h, fname)
			}

			results[i] = Result{
				Input:    fname,
				Hash:     hashStr,
				Error:    err,
				IsFile:   true,
				Duration: time.Since(start),
			}

			// Try to print results in order
			for {
				c := atomic.LoadInt32(&currFile)
				if c == i {
					printMu.Lock()
					// Print all consecutive completed results
					for int(c) < len(files) && results[c].Input != "" && !printed[c] {
						onResult(results[c])
						printed[c] = true
						c++
						atomic.StoreInt32(&currFile, c)
					}
					printMu.Unlock()
					break
				}
				// Check if our result is ready to be printed
				if c > i {
					break
				}
				time.Sleep(time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// Print any remaining results
	printMu.Lock()
	for i, r := range results {
		if !printed[i] {
			onResult(r)
		}
	}
	printMu.Unlock()
}

// computeFileHash streams a file through the hash
func computeFileHash(h hash.Hash, fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// hashPassword handles password-specific hashing algorithms
func hashPassword(password string, opts Options) (string, error) {
	switch opts.Algorithm {
	case "bcrypt":
		hash, err := bcrypt.GenerateFromPassword([]byte(password), opts.BcryptCost)
		if err != nil {
			return "", err
		}
		return string(hash), nil
	case "argon2id":
		// Generate a deterministic salt from the input for reproducible hashes
		// In production, you'd want a random salt stored with the hash
		salt := []byte("hashctl-argon2id-salt")
		hash := argon2.IDKey([]byte(password), salt, opts.Argon2Time, opts.Argon2Memory, opts.Argon2Lanes, opts.Argon2KeyLen)
		return hex.EncodeToString(hash), nil
	default:
		return "", fmt.Errorf("unsupported password hash: %s", opts.Algorithm)
	}
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}


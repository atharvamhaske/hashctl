// hashctl - A production-grade hashing CLI tool
// Inspired by CNCF-style infrastructure tools like kubectl, helm, and cosign
package main

import (
	"os"

	"github.com/atharvamhaske/hashctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

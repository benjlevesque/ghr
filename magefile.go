//+build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/sh"
)

// Build builds the code and generate a "ghr" executable
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.RunV("go", "build", "-o", "ghr", "-ldflags="+flags(), ".")
}

func flags() string {
	timestamp := time.Now().Format("2006-01-02")
	tag, _ := sh.Output("git", "describe", "--tags")
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	if tag == "" && hash == "" {
		tag = "dev"
	} else if tag == "" {
		tag = hash
	}
	return fmt.Sprintf(`-X "github.com/benjlevesque/ghr/pkg/cmd.timestamp=%s" -X "github.com/benjlevesque/ghr/pkg/cmd.gitTag=%s"`, timestamp, tag)
}

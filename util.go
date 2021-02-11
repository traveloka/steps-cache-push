package main

import (
    "os"
	"fmt"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"

	"time"
)

// Zip and split into few parts based on the size.
// Size can be 2g 20m. g mean gigabyte m mean megabyte
func Split(name, size, cacheArchivePath string) error {
    startTime := time.Now()

    cmd := command.New("split", "-b", size, cacheArchivePath, name)

    cmd.SetStdout(os.Stdout)
    cmd.SetStderr(os.Stderr)
    log.Debugf("$ " + cmd.PrintableCommandArgs())
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to split archive, error: %s", err)
    }

    log.Donef("Split done in %s\n", time.Since(startTime))

	return nil
}
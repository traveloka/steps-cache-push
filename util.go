package main

import (
    "os"
	"fmt"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"

	"time"
)

// Zip and split into few parts based on the size.
// Size 20m. m mean megabyte
func Split(name, size, cacheArchivePath string) error {
    startTime := time.Now()

    fi, err := os.Stat(cacheArchivePath)
	if err != nil {
		return fmt.Errorf("failed to get file info (%s): %s", cacheArchivePath, err)
	}
	sizeInBytes := fi.Size()
	log.Printf("Splitting archive with size: %d bytes / %f MB to %s", sizeInBytes, (float64(sizeInBytes) / 1024.0 / 1024.0), size)

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

func StoreCacheURL(url string) error {
    cmd := command.New("bash", "-c", "echo " + url + " > /Users/vagrant/cache.txt")

    cmd.SetStdout(os.Stdout)
    cmd.SetStderr(os.Stderr)
    log.Debugf("$ " + cmd.PrintableCommandArgs())
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to split archive, error: %s", err)
    }

    log.Donef("echo url success...\n")

    return nil
}
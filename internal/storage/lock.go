// Package storage provides persistence functionality for tasks
// in various formats including JSON and CSV.
package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/ZeRg0912/logger"
)

const (
	lockTimeout = 5 * time.Second
	lockRetry   = 100 * time.Millisecond
)

// FileLock represents a file lock for concurrent access protection.
type FileLock struct {
	lockFile *os.File
	path     string
}

// AcquireLock acquires an exclusive lock on a file.
// Returns an error if the lock cannot be acquired within the timeout.
func AcquireLock(path string) (*FileLock, error) {
	lockPath := path + ".lock"
	start := time.Now()

	for {
		file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			lock := &FileLock{
				lockFile: file,
				path:     lockPath,
			}
			logger.Debug("Acquired lock for %s", path)
			return lock, nil
		}

		if time.Since(start) > lockTimeout {
			return nil, fmt.Errorf("cannot acquire lock for %s: timeout after %v", path, lockTimeout)
		}

		time.Sleep(lockRetry)
	}
}

// Release releases the file lock.
func (fl *FileLock) Release() error {
	if fl.lockFile != nil {
		fl.lockFile.Close()
	}
	if err := os.Remove(fl.path); err != nil && !os.IsNotExist(err) {
		logger.Warn("Failed to remove lock file %s: %v", fl.path, err)
		return fmt.Errorf("cannot release lock: %w", err)
	}
	logger.Debug("Released lock for %s", fl.path)
	return nil
}


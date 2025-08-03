//go:build !windows

package main

import (
	"os"
	"syscall"
)

// mmapFile maps the given file into memory using Unix-specific syscalls.
func mmapFile(file *os.File) ([]byte, func(), error) {
	fi, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	fileSize := fi.Size()

	// Memory-map the file.
	data, err := syscall.Mmap(int(file.Fd()), 0, int(fileSize), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, nil, err
	}

	// Return a cleanup function that the caller should defer.
	cleanup := func() {
		syscall.Munmap(data)
	}

	return data, cleanup, nil
}

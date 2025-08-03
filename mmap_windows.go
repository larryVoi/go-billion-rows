//go:build windows

package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

// mmapFile maps the given file into memory using Windows-specific syscalls.
func mmapFile(file *os.File) ([]byte, func(), error) {
	fi, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	fileSize := fi.Size()

	// Create a file mapping handle.
	h, err := syscall.CreateFileMapping(syscall.Handle(file.Fd()), nil, syscall.PAGE_READONLY, 0, 0, nil)
	if err != nil {
		log.Fatalf("Error creating file mapping: %v", err)
	}
	// Note: The handle h from CreateFileMapping must be closed, but we can't
	// do it here because the mapping depends on it. In a real-world complex app,
	// you'd wrap this in a struct to manage the lifetime of the handle and the mapping.
	// For this single-file program, letting the OS clean up on exit is acceptable.

	// Map the file into the process's address space.
	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		log.Fatalf("Error mapping view of file: %v", err)
	}

	// It's crucial to unmap the view when we are done.
	// We return a cleanup function to be called by the main function.
	cleanup := func() {
		syscall.UnmapViewOfFile(addr)
		syscall.CloseHandle(h)
	}

	// Convert the mapped memory address to a byte slice.
	data := unsafe.Slice((*byte)(unsafe.Pointer(addr)), fileSize)

	return data, cleanup, nil
}

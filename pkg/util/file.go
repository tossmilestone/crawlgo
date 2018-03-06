package util

import (
	"io"
	"os"
)

// MkdirAllFunc is a function that makes directory with all in path.
type MkdirAllFunc func(path string, perm os.FileMode) error

// StatFunc is a function that get the info of a file.
type StatFunc func(name string) (os.FileInfo, error)

// CreateFunc is a function that create a file and return a io.Writer.
type CreateFunc func(name string) (io.Writer, error)

// DefaultMkdirAll is a default "MkdirAllFunc" implementation.
func DefaultMkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// DefaultStatFunc is a default "StatFunc" implementation.
func DefaultStatFunc(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// DefaultCreateFunc is a default "CreateFunc" implementation.
func DefaultCreateFunc(name string) (io.Writer, error) {
	return os.Create(name)
}

// MkdirAll is a MkdirAllFunc function.
var MkdirAll = DefaultMkdirAll

// Stat is a StatFunc function.
var Stat = DefaultStatFunc

// Create is a CreateFunc function.
var Create = DefaultCreateFunc

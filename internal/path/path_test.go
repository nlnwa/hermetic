package path

import (
	"os"
	"testing"
)

func TestIsDirectoryWithFile(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	defer os.Remove(file.Name())

	isDir, err := IsDirectory(file.Name())
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if isDir {
		t.Errorf("Expected '%s' to be a file", file.Name())
	}
}

func TestIsDirectoryWithDirectory(t *testing.T) {
	directory, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	defer os.Remove(directory)

	isDir, err := IsDirectory(directory)
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if !isDir {
		t.Errorf("Expected '%s' to be a directory", directory)
	}

}

func TestIsFileWithFile(t *testing.T) {
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	defer os.Remove(file.Name())

	isFile, err := IsFile(file.Name())
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if !isFile {
		t.Errorf("Expected '%s' to be a file", file.Name())
	}
}

func TestIsFileWithDirectory(t *testing.T) {
	directory, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	defer os.Remove(directory)

	isFile, err := IsFile(directory)
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if isFile {
		t.Errorf("Expected '%s' to be a directory", directory)
	}

}

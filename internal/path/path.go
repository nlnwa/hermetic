package path

import (
	"fmt"
	"os"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to get file info for '%s', original error: '%w'", path, err)
	}

	return fileInfo.IsDir(), err
}

func IsFile(path string) (bool, error) {
	isDir, err := IsDirectory(path)
	if err != nil {
		return false, fmt.Errorf("failed to check if '%s' is a directory, original error: '%w'", path, err)
	}
	return !isDir, err
}

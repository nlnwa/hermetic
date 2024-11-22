package path

import (
	"os"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func IsFile(path string) (bool, error) {
	isDir, err := IsDirectory(path)
	if err != nil {
		return false, err
	}

	return !isDir, nil
}

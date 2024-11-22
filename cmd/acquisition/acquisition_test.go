package acquisition

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidate(t *testing.T) {
	workingDirectory, err := os.MkdirTemp("", "working-dir")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	defer os.RemoveAll(workingDirectory)

	acquisitionYamlFile, err := os.Create(filepath.Join(workingDirectory, "acquisition.yaml"))
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	_, err = os.Create(filepath.Join(workingDirectory, "checksum.md5"))
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	_, err = os.Create(filepath.Join(workingDirectory, "checksum_transferred.md5"))
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	_, err = os.Create(filepath.Join(workingDirectory, "dummy.txt"))
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	_, err = acquisitionYamlFile.WriteString(`
__acquisition_version__: "0.2.0"
archive-unit:
    name: "dummy-acquisition"
    type: "acquisition"
    creator: "john doe"
    description: "Lorem ipsum"
    copyright-clearance: "Lorem ipsum"
    access-considerations: "Lorem ipsum"
    deposit:
        depositor: "jane doe"
        date: "2024-01-03T09:29:16+00:00"
        acquisition-purpose: "Lorem ipsum"
    handling:
        author: "jake doe"

files:
    - name: "acquisition.yaml"
      description: "This file"
      format: "yaml"
      path: "acquisition.yaml"
    - name: "checksum.md5"
      description: "Checksum file for this acquisition, generated with command 'find * -type f -print0 | sort -z | xargs -0 md5sum -b > /tmp/checksum.md5 && mv /tmp/checksum.md5 .'"
      format: "md5"
      path: "checksum.md5"
    - name: "checksum_transferred.md5"
      description: "Checksum file for packaged files"
      format: "md5"
      path: "checksum_transferred.md5"
    - name: "dummy"
      description: "Dummy file"
      format: "plain"
      path: "dummy.txt"
`)

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	dataModel, err := deserializeYamlFile(acquisitionYamlFile.Name())
	if err != nil {
		t.Errorf("Expected no error, got '%s' %s", err, acquisitionYamlFile.Name())
	}

	if err := validate(workingDirectory, dataModel); err != nil {
		t.Errorf("Expected no error, got '%s' %s", err, acquisitionYamlFile.Name())
	}

}

package fs

import (
	"os"
	"path/filepath"
)

var _ Filesystem = &MockFilesystem{}

// A MockFilesystem is a filesystem held in-memory.
type MockFilesystem struct {
	Files map[string][]byte
}

// NewMock constructs a new mock filesystem.
func NewMock(files map[string][]byte) *MockFilesystem {
	return &MockFilesystem{Files: files}
}

// List iterates through all mock files and returns mock files with the
// directory as a prefix that do not contain any subsequent sections.
func (fs *MockFilesystem) List(dir string) ([]string, error) {
	var ls []string
	for f := range fs.Files {
		if filepath.Dir(f) == dir {
			ls = append(ls, filepath.Base(f))
		}
	}
	return ls, nil
}

// Read returns the contents of a mock file. If the file is missing, it returns
// an error that returns true when passed as an argument to os.IsNotExist.
func (fs *MockFilesystem) Read(filename string) ([]byte, error) {
	data, ok := fs.Files[filename]
	if !ok {
		return nil, os.ErrNotExist
	}
	return data, nil
}

// Delete deletes a mock file.
func (fs *MockFilesystem) Delete(filename string) error {
	delete(fs.Files, filename)
	return nil
}

package fs

import (
	"io/ioutil"
	"os"
)

var _ Filesystem = &OSFilesystem{}

// An OSFilesystem provides a Filesystem that interacts with the process's local
// underlying filesystem.
type OSFilesystem struct{}

// New constructs a new OSFilesystem instance.
func New() *OSFilesystem {
	return &OSFilesystem{}
}

// List delegates to ioutil.ReadDir.
func (*OSFilesystem) List(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

// Read delegates to ioutil.Read.
func (*OSFilesystem) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

// Delete delegates to os.Remove.
func (*OSFilesystem) Delete(filename string) error {
	return os.Remove(filename)
}

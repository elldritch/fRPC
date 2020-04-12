package fs

import (
	"fmt"
	"io"
)

var _ Filesystem = &DevFS{}

// A DevFS is an OSFilesystem that never deletes things, instead logging a
// message when a deletion would occur.
type DevFS struct {
	underlying OSFilesystem
	writer     io.Writer
}

// NewDevFS constructs a new DevFS instance.
func NewDevFS(writer io.Writer) *DevFS {
	return &DevFS{underlying: OSFilesystem{}, writer: writer}
}

// List delegates to OSFilesystem.List.
func (fs *DevFS) List(dir string) ([]string, error) {
	return fs.underlying.List(dir)
}

// Read delegates to OSFilesystem.Read.
func (fs *DevFS) Read(filename string) ([]byte, error) {
	return fs.underlying.Read(filename)
}

// Delete logs a message.
func (fs *DevFS) Delete(filename string) error {
	fmt.Fprintf(fs.writer, "Deleting file: %#v\n", filename)
	// info, err := os.Stat(filename)
	// if err != nil {
	// 	fmt.Fprintf(fs.writer, "Actual underlying file returned error on os.Stat: %#v\n", err)
	// } else {
	// 	fmt.Fprintf(fs.writer, "Actual underlying file: %#v\n", info)
	// }

	return nil
}

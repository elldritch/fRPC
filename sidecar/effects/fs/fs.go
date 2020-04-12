// Package fs provides an effect abstraction for using a filesystem.
package fs

// A Filesystem provides functionality for using a filesystem.
type Filesystem interface {
	List(dir string) ([]string, error)
	Read(filename string) ([]byte, error)
	Delete(filename string) error
}

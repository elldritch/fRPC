package fs_test

import (
	"testing"

	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/stretchr/testify/assert"
)

func TestOSFilesystemList(t *testing.T) {
	fs := fs.New()
	files, err := fs.List("./testdata")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, files)
}

func TestOSFilesystemRead(t *testing.T) {
	fs := fs.New()

	data, err := fs.Read("./testdata/a")
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello, world!\n"), data)

	data, err = fs.Read("./testdata/b")
	assert.NoError(t, err)
	assert.Equal(t, []byte("another file\n"), data)
}

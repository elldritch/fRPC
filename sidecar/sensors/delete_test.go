package sensors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/liftM/fRPC/sidecar/sensors"
)

func TestExpired(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json": nil,
		"2.json": nil,
		"3.json": nil,
	})

	expired, err := sensors.Expired(fs, ".", 1)
	assert.NoError(t, err)
	assert.Len(t, expired, 2)
	assert.Contains(t, expired, "1.json")
	assert.Contains(t, expired, "2.json")
}

func TestExpiredDir(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"dir/1.json": nil,
		"dir/2.json": nil,
		"dir/3.json": nil,
	})

	expired, err := sensors.Expired(fs, "dir", 1)
	assert.NoError(t, err)
	assert.Len(t, expired, 2)
	assert.Contains(t, expired, "dir/1.json")
	assert.Contains(t, expired, "dir/2.json")
}

func TestExpiredEmpty(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{})

	expired, err := sensors.Expired(fs, ".", 1)
	assert.NoError(t, err)
	assert.Len(t, expired, 0)
}

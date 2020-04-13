package sensors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/liftM/fRPC/sidecar/sensors"
)

func TestFRPCSensorDeleteExpired(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json":   nil,
		"100.json": nil,
	})

	sensor := sensors.New(sensors.Config{
		Filesystem: fs,
		Dir:        ".",
		TTL:        sensors.ToDuration(1),
	})

	err := sensor.DeleteExpired()
	assert.NoError(t, err)
	assert.Len(t, fs.Files, 1)
	assert.NotContains(t, fs.Files, "1.json")
	assert.Contains(t, fs.Files, "100.json")
}

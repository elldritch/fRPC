package sensors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liftM/fRPC/sidecar/effects/fs"
	"github.com/liftM/fRPC/sidecar/sensors"
)

func TestUnmarshal(t *testing.T) {
	sample, err := sensors.Unmarshal([]byte(`{"tick":7996904,"values":[{"network_id":4,"signals":[{"signal":{"type":"item","name":"copper-plate"},"count":1}]}]}`))
	assert.NoError(t, err)
	assert.Equal(t, sensors.Tick(7996904), sample.Tick)
	assert.Equal(t, sensors.Count(1), sample.Readings[4]["copper-plate"])
}

func TestUnmarshalMultipleReadingsSameNetwork(t *testing.T) {
	sample, err := sensors.Unmarshal([]byte(`{"tick":8006492,"values":[{"network_id":4,"signals":[{"signal":{"type":"item","name":"copper-ore"},"count":4}]},{"network_id":4,"signals":[{"signal":{"type":"item","name":"copper-ore"},"count":4}]}]}`))
	assert.NoError(t, err)
	assert.Equal(t, sensors.Tick(8006492), sample.Tick)
	assert.Equal(t, sensors.Count(4), sample.Readings[4]["copper-ore"])
}

func TestReadFilesAscendingOrderByContentsNotFilename(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json": []byte(`{"tick": 3, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 4}]}]}`),
		"2.json": []byte(`{"tick": 2, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 5}]}]}`),
		"3.json": []byte(`{"tick": 1, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 6}]}]}`),
	})

	samples, err := sensors.ReadFiles(fs, []string{"1.json", "2.json", "3.json"})
	assert.NoError(t, err)
	assert.Len(t, samples, 3)
	assert.Equal(t, sensors.Count(6), samples[0].Readings[1]["test"])
	assert.Equal(t, sensors.Count(5), samples[1].Readings[1]["test"])
	assert.Equal(t, sensors.Count(4), samples[2].Readings[1]["test"])
}

func TestReadFilesDir(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"dir/1.json": []byte(`{"tick": 3, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 4}]}]}`),
		"dir/2.json": []byte(`{"tick": 2, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 5}]}]}`),
		"dir/3.json": []byte(`{"tick": 1, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 6}]}]}`),
	})

	samples, err := sensors.ReadFiles(fs, []string{"dir/1.json", "dir/2.json", "dir/3.json"})
	assert.NoError(t, err)
	assert.Len(t, samples, 3)
	assert.Equal(t, sensors.Count(6), samples[0].Readings[1]["test"])
	assert.Equal(t, sensors.Count(5), samples[1].Readings[1]["test"])
	assert.Equal(t, sensors.Count(4), samples[2].Readings[1]["test"])
}

func TestSinceTick(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json": []byte(`{"tick": 1, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 4}]}]}`),
		"2.json": []byte(`{"tick": 2, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 5}]}]}`),
		"3.json": []byte(`{"tick": 3, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 6}]}]}`),
	})

	samples, err := sensors.Since(fs, ".", 2, 10)
	assert.NoError(t, err)
	assert.Len(t, samples, 2)
	assert.Equal(t, sensors.Count(5), samples[0].Readings[1]["test"])
	assert.Equal(t, sensors.Count(6), samples[1].Readings[1]["test"])
}

func TestSinceCount(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json": []byte(`{"tick": 1, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 4}]}]}`),
		"2.json": []byte(`{"tick": 2, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 5}]}]}`),
		"3.json": []byte(`{"tick": 3, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 6}]}]}`),
	})

	samples, err := sensors.Since(fs, ".", 1, 2)
	assert.NoError(t, err)
	assert.Len(t, samples, 2)
	assert.Equal(t, sensors.Count(4), samples[0].Readings[1]["test"])
	assert.Equal(t, sensors.Count(5), samples[1].Readings[1]["test"])
}

func TestSinceDir(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"dir/1.json": []byte(`{"tick": 1, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 4}]}]}`),
		"dir/2.json": []byte(`{"tick": 2, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 5}]}]}`),
		"dir/3.json": []byte(`{"tick": 3, "values": [{"network_id": 1, "signals": [{"signal": {"type": "fake", "name": "test"}, "count": 6}]}]}`),
	})

	samples, err := sensors.Since(fs, "dir", 1, 2)
	assert.NoError(t, err)
	assert.Len(t, samples, 2)
	assert.Equal(t, sensors.Count(4), samples[0].Readings[1]["test"])
	assert.Equal(t, sensors.Count(5), samples[1].Readings[1]["test"])
}

func TestLatestTick(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"1.json": nil,
		"2.json": nil,
		"3.json": nil,
	})

	tick, err := sensors.LatestTick(fs, ".")
	assert.NoError(t, err)
	assert.Equal(t, sensors.Tick(3), tick)
}

func TestLatestTickDir(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{
		"dir/1.json": nil,
		"dir/2.json": nil,
		"dir/3.json": nil,
	})

	tick, err := sensors.LatestTick(fs, "dir")
	assert.NoError(t, err)
	assert.Equal(t, sensors.Tick(3), tick)
}

func TestLatestTickEmpty(t *testing.T) {
	fs := fs.NewMock(map[string][]byte{})

	tick, err := sensors.LatestTick(fs, ".")
	assert.NoError(t, err)
	assert.Equal(t, sensors.Tick(0), tick)
}

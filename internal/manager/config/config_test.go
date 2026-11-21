package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func TestConfig_GetAllPluginConfiguration(t *testing.T) {
	i := InitializeEmpty()

	assert.Equal(t, i.GetAllPluginConfiguration(), map[string]map[string]interface{}{})

	i.SetPluginConfiguration("plugin1", map[string]interface{}{"key1": "value1"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
	}, i.GetAllPluginConfiguration())

	i.SetPluginConfiguration("plugin2", map[string]interface{}{"key2": "value2"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
		"plugin2": {"key2": "value2"},
	}, i.GetAllPluginConfiguration())

	// ensure SetPluginConfiguration overwrites existing configuration
	i.SetPluginConfiguration("plugin2", map[string]interface{}{"key3": "value3"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
		"plugin2": {"key3": "value3"},
	}, i.GetAllPluginConfiguration())
}

func TestConfig_initialisePublicWhitelist(t *testing.T) {
	i := InitializeEmpty()
	i.SetInterface(PublicWhitelist, []string{"203.0.113.45", "198.51.100.0/24", "2001:db8::/32"})

	assert.NoError(t, i.initialisePublicWhitelist())
	nets, addrs := i.GetPublicWhitelist()
	assert.Len(t, nets, 2)
	assert.Len(t, addrs, 1)
	assert.Equal(t, "203.0.113.45", addrs[0].String())

	for _, bad := range []string{"/0", "not-an-ip", "*"} {
		i := InitializeEmpty()
		i.SetInterface(PublicWhitelist, []string{bad})
		assert.Error(t, i.initialisePublicWhitelist(), bad)
	}
}

func TestConfig_GetGeneratedVideoCodec(t *testing.T) {
	i := InitializeEmpty()

	assert.Equal(t, models.GeneratedVideoCodecH264, i.GetPreviewCodec(), "unset defaults to H264")

	i.SetString(PreviewCodec, "HEVC")
	assert.Equal(t, models.GeneratedVideoCodecHevc, i.GetPreviewCodec())

	// a hand-edited config.yml must not leak an unknown codec past the whitelist
	i.SetString(TranscodeCodec, "vp9")
	assert.Equal(t, models.GeneratedVideoCodecH264, i.GetTranscodeCodec())
}

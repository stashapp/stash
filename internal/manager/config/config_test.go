package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetPassword(t *testing.T) {
	i := InitializeEmpty()

	pw := strings.Repeat("a", maxPasswordBytes)
	assert.NoError(t, i.SetPassword(pw))
	assert.NotEmpty(t, i.GetPasswordHash())

	// error must not overwrite the existing hash
	before := i.GetPasswordHash()
	assert.Error(t, i.SetPassword(pw+"a"))
	assert.Equal(t, before, i.GetPasswordHash())

	// 36 two-byte runes = 72 bytes ok; 37 = 74 bytes rejected
	assert.NoError(t, i.SetPassword(strings.Repeat("é", 36)))
	assert.Error(t, i.SetPassword(strings.Repeat("é", 37)))

	assert.NoError(t, i.SetPassword(""))
	assert.Empty(t, i.GetPasswordHash())
}

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

func TestConfig_GetTranscodeSkipCodecs(t *testing.T) {
	i := InitializeEmpty()

	assert.Equal(t, []string{"h264"}, i.GetTranscodeSkipCodecs())

	i.SetInterface(TranscodeSkipCodecs, []string{" HEVC ", "", "av1"})
	assert.Equal(t, []string{"hevc", "av1"}, i.GetTranscodeSkipCodecs())

	i.SetInterface(TranscodeSkipCodecs, []string{" "})
	assert.Equal(t, []string{"h264"}, i.GetTranscodeSkipCodecs())
}

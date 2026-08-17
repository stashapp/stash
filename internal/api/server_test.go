package api

import (
	"testing"

	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestCspConnectSrcFromSettings(t *testing.T) {
	for _, tt := range []struct {
		name        string
		settings    map[string]interface{}
		valid       []string
		skippedKeys []string
	}{
		{
			name:     "no csp keys",
			settings: map[string]interface{}{"foo": "bar", "other_setting": "https://x.com"},
			valid:    nil,
		},
		{
			name:     "valid https url",
			settings: map[string]interface{}{"csp_x": "https://api.example.com"},
			valid:    []string{"https://api.example.com"},
		},
		{
			name:     "valid http url",
			settings: map[string]interface{}{"csp_x": "http://localhost:7860"},
			valid:    []string{"http://localhost:7860"},
		},
		{
			name:        "disallowed scheme",
			settings:    map[string]interface{}{"csp_x": "ftp://example.com"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "not a url",
			settings:    map[string]interface{}{"csp_x": "not a url"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "non string value",
			settings:    map[string]interface{}{"csp_x": 123},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "wildcard host",
			settings:    map[string]interface{}{"csp_x": "http://*:7860"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "csp directive breakout via semicolon in path",
			settings:    map[string]interface{}{"csp_x": "https://evil.com/; script-src 'none'"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "whitespace in path",
			settings:    map[string]interface{}{"csp_x": "https://evil.com/a b"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "comma in url",
			settings:    map[string]interface{}{"csp_x": "https://evil.com/a,b"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "degenerate host port only",
			settings:    map[string]interface{}{"csp_x": "https://:7860"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "userinfo in url",
			settings:    map[string]interface{}{"csp_x": "https://user@attacker.com"},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "empty string",
			settings:    map[string]interface{}{"csp_x": ""},
			skippedKeys: []string{"csp_x"},
		},
		{
			name:        "mixed valid and invalid",
			settings:    map[string]interface{}{"csp_a": "https://api.example.com", "csp_b": "javascript:alert(1)", "csp_c": "http://localhost:7860", "other": "ignored"},
			valid:       []string{"https://api.example.com", "http://localhost:7860"},
			skippedKeys: []string{"csp_b"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			valid, skippedKeys := cspConnectSrcFromSettings(tt.settings)
			assert.ElementsMatch(t, tt.valid, valid)
			assert.ElementsMatch(t, tt.skippedKeys, skippedKeys)
		})
	}
}

func TestSetPageSecurityHeaders_CSPSettingsOptIn(t *testing.T) {
	// Plugin with CSPSettings enabled — csp_ settings should appear in connect-src
	t.Run("opt-in enabled", func(t *testing.T) {
		plugins := []*plugin.Plugin{
			{
				Enabled: true,
				UI: plugin.PluginUI{
					CSPSettings: true,
				},
			},
		}

		assert.True(t, plugins[0].UI.CSPSettings)
	})

	// Plugin without CSPSettings — csp_ settings should NOT be scanned
	t.Run("opt-in disabled", func(t *testing.T) {
		plugins := []*plugin.Plugin{
			{
				Enabled: true,
				UI: plugin.PluginUI{
					CSPSettings: false,
				},
			},
		}

		assert.False(t, plugins[0].UI.CSPSettings)
	})
}

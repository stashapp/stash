package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestCspConnectSrcFromSettings(t *testing.T) {
	for _, tt := range []struct {
		name     string
		settings map[string]interface{}
		valid    []string
		skipped  map[string]string
	}{
		{
			name:     "no csp keys",
			settings: map[string]interface{}{"foo": "bar", "other_setting": "https://x.com"},
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
			name:     "disallowed scheme",
			settings: map[string]interface{}{"csp_x": "ftp://example.com"},
			skipped:  map[string]string{"csp_x": "ftp://example.com"},
		},
		{
			name:     "not a url",
			settings: map[string]interface{}{"csp_x": "not a url"},
			skipped:  map[string]string{"csp_x": "not a url"},
		},
		{
			name:     "non string value",
			settings: map[string]interface{}{"csp_x": 123},
			skipped:  map[string]string{"csp_x": "123"},
		},
		{
			name:     "wildcard host",
			settings: map[string]interface{}{"csp_x": "http://*:7860"},
			skipped:  map[string]string{"csp_x": "http://*:7860"},
		},
		{
			name:     "csp directive breakout via semicolon in path",
			settings: map[string]interface{}{"csp_x": "https://evil.com/; script-src 'none'"},
			skipped:  map[string]string{"csp_x": "https://evil.com/; script-src 'none'"},
		},
		{
			name:     "whitespace in path",
			settings: map[string]interface{}{"csp_x": "https://evil.com/a b"},
			skipped:  map[string]string{"csp_x": "https://evil.com/a b"},
		},
		{
			name:     "comma in url",
			settings: map[string]interface{}{"csp_x": "https://evil.com/a,b"},
			skipped:  map[string]string{"csp_x": "https://evil.com/a,b"},
		},
		{
			name:     "degenerate host port only",
			settings: map[string]interface{}{"csp_x": "https://:7860"},
			skipped:  map[string]string{"csp_x": "https://:7860"},
		},
		{
			name:     "userinfo in url",
			settings: map[string]interface{}{"csp_x": "https://user@attacker.com"},
			skipped:  map[string]string{"csp_x": "https://user@attacker.com"},
		},
		{
			name:     "empty string",
			settings: map[string]interface{}{"csp_x": ""},
			skipped:  map[string]string{"csp_x": ""},
		},
		{
			name:     "mixed valid and invalid",
			settings: map[string]interface{}{"csp_a": "https://api.example.com", "csp_b": "javascript:alert(1)", "csp_c": "http://localhost:7860", "other": "ignored"},
			valid:    []string{"http://localhost:7860", "https://api.example.com"},
			skipped:  map[string]string{"csp_b": "javascript:alert(1)"},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			valid, skipped := cspConnectSrcFromSettings(tt.settings)
			// valid is sorted, so the emitted header is stable between requests
			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.skipped, skipped)
		})
	}
}

// connectSrc returns the connect-src directive of the CSP header emitted for a
// page request, given the supplied plugins and stored plugin configuration.
func connectSrc(t *testing.T, plugins []*plugin.Plugin, pluginConfig map[string]interface{}) string {
	t.Helper()

	c := config.InitializeEmpty()
	for _, p := range plugins {
		c.SetPluginConfiguration(p.ID, pluginConfig)
	}

	w := httptest.NewRecorder()
	setPageSecurityHeaders(w, httptest.NewRequest(http.MethodGet, "/", nil), plugins)

	return w.Header().Get("Content-Security-Policy")
}

func TestSetPageSecurityHeadersCSPSettings(t *testing.T) {
	settings := map[string]interface{}{
		"csp_endpoint": "https://api.example.com",
		"csp_bad":      "http://*:7860",
		"apiKey":       "secret",
	}

	pluginWith := func(cspSettings bool) []*plugin.Plugin {
		return []*plugin.Plugin{{
			ID:      "test-plugin",
			Enabled: true,
			UI:      plugin.PluginUI{CSPSettings: cspSettings},
		}}
	}

	t.Run("opted in", func(t *testing.T) {
		csp := connectSrc(t, pluginWith(true), settings)
		assert.Contains(t, csp, "https://api.example.com")
		assert.NotContains(t, csp, "http://*:7860")
		assert.NotContains(t, csp, "secret")
	})

	t.Run("not opted in", func(t *testing.T) {
		csp := connectSrc(t, pluginWith(false), settings)
		assert.NotContains(t, csp, "https://api.example.com")
	})

	t.Run("disabled plugin", func(t *testing.T) {
		plugins := pluginWith(true)
		plugins[0].Enabled = false
		assert.NotContains(t, connectSrc(t, plugins, settings), "https://api.example.com")
	})
}

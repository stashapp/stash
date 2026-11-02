package api

import (
	"context"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
)

func TestIsAdminRequest(t *testing.T) {
	withRoles := func(roles ...models.RoleEnum) context.Context {
		return session.SetCurrentUser(context.Background(), models.User{ID: 1, Username: "u", Roles: roles})
	}

	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		{"no user", context.Background(), true},
		{"admin", withRoles(models.RoleEnumAdmin), true},
		{"read", withRoles(models.RoleEnumRead), false},
		{"modify", withRoles(models.RoleEnumModify), false},
		{"no roles", withRoles(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAdminRequest(tt.ctx); got != tt.want {
				t.Errorf("isAdminRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedactConfigResult(t *testing.T) {
	str := func(s string) *string { return &s }

	box := &models.StashBox{Endpoint: "https://stashdb.org/graphql", APIKey: "secret", Name: "StashDB"}
	ret := &ConfigResult{
		General: &ConfigGeneralResult{
			Stashes:                      []*config.StashConfig{{Path: "/library"}},
			DatabasePath:                 "/data/stash.sqlite",
			DeleteTrashPath:              "/trash",
			FfmpegPath:                   "/usr/bin/ffmpeg",
			PythonPath:                   "/usr/bin/python",
			LogFile:                      str("/data/stash.log"),
			CustomPerformerImageLocation: str("/images"),
			TranscodeInputArgs:           []string{"-hwaccel", "auto"},
			StashBoxes:                   []*models.StashBox{box},
		},
		Scraping: &ConfigScrapingResult{
			ScraperCDPPath: str("http://user:pass@chrome:9222"),
		},
	}

	redactConfigResult(ret)

	g := ret.General
	if g.DatabasePath != "" || g.FfmpegPath != "" || g.PythonPath != "" {
		t.Errorf("server paths not redacted: %q %q %q", g.DatabasePath, g.FfmpegPath, g.PythonPath)
	}
	if g.LogFile != nil || g.CustomPerformerImageLocation != nil {
		t.Error("optional paths not redacted")
	}
	if len(g.TranscodeInputArgs) != 0 {
		t.Errorf("transcode args not redacted: %v", g.TranscodeInputArgs)
	}
	if ret.Scraping.ScraperCDPPath != nil {
		t.Error("scraper CDP path not redacted")
	}

	// values used by the UI for non-administrators are kept
	if len(g.Stashes) != 1 || g.Stashes[0].Path != "/library" {
		t.Error("library paths should be kept")
	}
	if g.DeleteTrashPath != "/trash" {
		t.Error("trash path should be kept")
	}

	if len(g.StashBoxes) != 1 {
		t.Fatalf("expected 1 stash box, got %d", len(g.StashBoxes))
	}
	got := g.StashBoxes[0]
	if got.APIKey != "" {
		t.Error("stash-box API key not redacted")
	}
	if got.Endpoint != box.Endpoint || got.Name != box.Name {
		t.Error("stash-box endpoint and name should be kept")
	}

	// the source stash box must not be modified, since it may be owned by the config
	if box.APIKey != "secret" {
		t.Error("redaction modified the source stash box")
	}
}

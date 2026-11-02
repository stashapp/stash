package user

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
)

func TestLegacyService_GetSingleUser(t *testing.T) {
	tests := []struct {
		name         string
		service      LegacyService
		wantUsername string
		wantNil      bool
	}{
		{"no credentials", LegacyService{}, legacyDefaultUsername, false},
		{"username without password", LegacyService{Username: "user"}, "user", false},
		{"credentials", LegacyService{Username: "user", PasswordHash: "hash"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := tt.service.GetSingleUser(context.Background())
			if err != nil {
				t.Fatalf("GetSingleUser() error = %v", err)
			}

			if tt.wantNil {
				if u != nil {
					t.Errorf("GetSingleUser() = %v, want nil", u)
				}
				return
			}

			if u == nil {
				t.Fatal("GetSingleUser() = nil, want user")
			}
			if u.Username != tt.wantUsername {
				t.Errorf("GetSingleUser().Username = %q, want %q", u.Username, tt.wantUsername)
			}
			if !u.Roles.HasRole(models.RoleEnumAdmin) {
				t.Errorf("GetSingleUser().Roles = %v, want admin", u.Roles)
			}
		})
	}
}

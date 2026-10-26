package models

import (
	"context"
	"testing"
)

func TestUserIDFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		wantID int
		wantOK bool
	}{
		{"not set", context.Background(), 0, false},
		{"guest user", WithUserID(context.Background(), 0), 0, false},
		{"stored user", WithUserID(context.Background(), 3), 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotOK := UserIDFromContext(tt.ctx)
			if gotID != tt.wantID || gotOK != tt.wantOK {
				t.Errorf("UserIDFromContext() = (%d, %v), want (%d, %v)", gotID, gotOK, tt.wantID, tt.wantOK)
			}
		})
	}
}

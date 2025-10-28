package models

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateIDs_ImpactedIDs(t *testing.T) {
	tests := []struct {
		name     string
		IDs      []int
		Mode     RelationshipUpdateMode
		existing []int
		want     []int
	}{
		{
			name:     "add",
			IDs:      []int{1, 2, 3},
			Mode:     RelationshipUpdateModeAdd,
			existing: []int{1, 2},
			want:     []int{3},
		},
		{
			name:     "remove",
			IDs:      []int{1, 2, 3},
			Mode:     RelationshipUpdateModeRemove,
			existing: []int{1, 2},
			want:     []int{1, 2},
		},
		{
			name:     "set",
			IDs:      []int{1, 2, 3},
			Mode:     RelationshipUpdateModeSet,
			existing: []int{1, 2},
			want:     []int{3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UpdateIDs{
				IDs:  tt.IDs,
				Mode: tt.Mode,
			}
			if got := u.ImpactedIDs(tt.existing); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateIDs.ImpactedIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyUpdate(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		mode     RelationshipUpdateMode
		existing []int
		want     []int
	}{
		{
			name:     "add",
			values:   []int{2, 3},
			mode:     RelationshipUpdateModeAdd,
			existing: []int{1, 2},
			want:     []int{1, 2, 3},
		},
		{
			name:     "remove",
			values:   []int{2, 3},
			mode:     RelationshipUpdateModeRemove,
			existing: []int{1, 2},
			want:     []int{1},
		},
		{
			name:     "set",
			values:   []int{1, 2, 3},
			mode:     RelationshipUpdateModeSet,
			existing: []int{1, 2},
			want:     []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyUpdate(tt.values, tt.mode, tt.existing)
			assert.Equal(t, tt.want, got)
		})
	}
}

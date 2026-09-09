package scene

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stashapp/stash/pkg/models"
)

func makeVideoFile(id models.FileID, duration float64) *models.VideoFile {
	return &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:   id,
			Path: "/test/file.mp4",
		},
		Duration: duration,
	}
}

func TestValidateFileRange(t *testing.T) {
	duration := 100.0
	f := makeVideoFile(1, duration)

	var (
		zero  = 0.0
		five  = 5.0
		fifty = 50.0
		hund  = 100.0
		over  = 150.0
	)

	tests := []struct {
		name      string
		start     *float64
		end       *float64
		wantErr   bool
		errSubstr string
	}{
		{"whole file", nil, nil, false, ""},
		{"open start", nil, &fifty, false, ""},
		{"open end", &fifty, nil, false, ""},
		{"both bounds", &five, &fifty, false, ""},
		{"start at zero", &zero, &fifty, false, ""},
		{"end at duration", &fifty, &hund, false, ""},
		{"start below zero", func() *float64 { v := -1.0; return &v }(), nil, true, "must not be negative"},
		{"start below zero", func() *float64 { v := -1.0; return &v }(), nil, true, "must not be negative"},
		{"zero end", &zero, &zero, true, "must be positive"},
		{"start after end", &fifty, &five, true, "must be before"},
		{"start equals end", &fifty, &fifty, true, "must be before"},
		{"start beyond duration", &over, nil, true, "end of file"},
		{"end beyond duration", nil, &over, true, "duration of file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileRange(f, tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errSubstr != "" {
					assert.Contains(t, err.Error(), tt.errSubstr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFileRange_UnprobedFile(t *testing.T) {
	// files with unknown duration must not be rejected outright
	f := makeVideoFile(1, 0)

	start := 10.0
	err := validateFileRange(f, &start, nil)
	assert.NoError(t, err)
}

func TestValidateFileRanges(t *testing.T) {
	f1 := makeVideoFile(1, 100)
	f2 := makeVideoFile(2, 200)
	files := []*models.VideoFile{f1, f2}

	fileID1 := models.FileID(1)
	fileID2 := models.FileID(2)

	var (
		ten  = 10.0
		twen = 20.0
	)

	primary := &fileID1

	t.Run("resolves nil file id to primary", func(t *testing.T) {
		ret, err := validateFileRanges([]models.FileID{1, 2}, primary, files, []FileRangeInput{
			{FileID: nil, StartTime: &ten, EndTime: &twen},
		})
		assert.NoError(t, err)
		assert.Len(t, ret, 1)
		r, ok := ret[fileID1]
		assert.True(t, ok)
		assert.Equal(t, &ten, r.StartTime)
	})

	t.Run("explicit file id", func(t *testing.T) {
		ret, err := validateFileRanges([]models.FileID{1, 2}, primary, files, []FileRangeInput{
			{FileID: &fileID2, StartTime: &ten},
		})
		assert.NoError(t, err)
		_, ok := ret[fileID2]
		assert.True(t, ok)
	})

	t.Run("no primary and nil file id", func(t *testing.T) {
		_, err := validateFileRanges(nil, nil, files, []FileRangeInput{
			{FileID: nil},
		})
		assert.ErrorIs(t, err, errNoPrimaryFile)
	})

	t.Run("duplicate range for same file", func(t *testing.T) {
		_, err := validateFileRanges([]models.FileID{1, 2}, primary, files, []FileRangeInput{
			{FileID: &fileID1, StartTime: &ten},
			{FileID: &fileID1, StartTime: &twen},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})

	t.Run("file not on scene", func(t *testing.T) {
		fileID3 := models.FileID(3)
		_, err := validateFileRanges([]models.FileID{1, 2}, primary, files, []FileRangeInput{
			{FileID: &fileID3},
		})
		assert.ErrorIs(t, err, ErrFileRangeNotFound)
	})

	t.Run("invalid bounds reported", func(t *testing.T) {
		_, err := validateFileRanges([]models.FileID{1, 2}, primary, files, []FileRangeInput{
			{FileID: &fileID1, StartTime: &twen, EndTime: &ten},
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be before")
	})
}

func TestSceneFileRange_Bounds(t *testing.T) {
	var (
		ten  = 10.0
		twen = 20.0
	)

	r := models.SceneFileRange{StartTime: &ten, EndTime: &twen}
	start, end := r.Bounds(100)
	assert.Equal(t, 10.0, start)
	assert.Equal(t, 20.0, end)
	assert.Equal(t, 10.0, r.Duration(100))

	r = models.SceneFileRange{}
	start, end = r.Bounds(100)
	assert.Equal(t, 0.0, start)
	assert.Equal(t, 100.0, end)
	assert.Equal(t, 100.0, r.Duration(100))

	r = models.SceneFileRange{StartTime: &ten}
	start, end = r.Bounds(100)
	assert.Equal(t, 10.0, start)
	assert.Equal(t, 100.0, end)

	r = models.SceneFileRange{EndTime: &twen}
	start, end = r.Bounds(100)
	assert.Equal(t, 0.0, start)
	assert.Equal(t, 20.0, end)
}

func TestFloatPtrEq(t *testing.T) {
	var (
		a = 1.5
		b = 1.5
		c = 2.5
	)

	assert.True(t, floatPtrEq(nil, nil))
	assert.True(t, floatPtrEq(&a, &b))
	assert.False(t, floatPtrEq(&a, &c))
	assert.False(t, floatPtrEq(&a, nil))
	assert.False(t, floatPtrEq(nil, &c))
}

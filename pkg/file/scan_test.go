package file

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
)

const (
	createFilePath    = "createFilePath"
	errCreateFilePath = "errCreateFilePath"
	updateFilePath    = "updateFilePath"
	errUpdateFilePath = "errUpdateFilePath"

	createFileID = iota + 1
)

func TestScanner_ApplyChanges(t *testing.T) {
	createFile := models.File{Path: createFilePath}
	errCreateFile := models.File{Path: errCreateFilePath}
	updateFile := models.File{Path: updateFilePath}
	errUpdateFile := models.File{Path: errUpdateFilePath}

	errCreate := errors.New("error creating file")
	errUpdate := errors.New("error updating file")

	rw := &mocks.FileReaderWriter{}
	rw.On("Create", createFile).Return(&models.File{ID: createFileID}, nil).Once()
	rw.On("Create", errCreateFile).Return(nil, errCreate).Once()
	rw.On("UpdateFull", updateFile).Return(&updateFile, nil).Once()
	rw.On("UpdateFull", errUpdateFile).Return(nil, errUpdate).Once()

	scanner := &Scanner{}

	tests := []struct {
		name    string
		scanned *Scanned
		wantErr bool
	}{
		{
			"create",
			&Scanned{
				New: &createFile,
			},
			false,
		},
		{
			"create error",
			&Scanned{
				New: &errCreateFile,
			},
			true,
		},
		{
			"update",
			&Scanned{
				Old: &updateFile,
				New: &updateFile,
			},
			false,
		},
		{
			"update error",
			&Scanned{
				Old: &updateFile,
				New: &errUpdateFile,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := scanner.ApplyChanges(rw, tt.scanned); (err != nil) != tt.wantErr {
				t.Errorf("Scanner.ApplyChanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	rw.AssertExpectations(t)
}

func TestScanned_FileUpdated(t *testing.T) {
	oldFile := &models.File{
		Path: createFilePath,
	}
	sameFile := &models.File{
		Path: createFilePath,
	}
	changedFile := &models.File{
		Path: updateFilePath,
	}

	tests := []struct {
		name string
		s    Scanned
		want bool
	}{
		{
			"old nil",
			Scanned{
				Old: nil,
				New: sameFile,
			},
			false,
		},
		{
			"new nil",
			Scanned{
				Old: oldFile,
				New: nil,
			},
			false,
		},
		{
			"same",
			Scanned{
				Old: oldFile,
				New: sameFile,
			},
			false,
		},
		{
			"different",
			Scanned{
				Old: oldFile,
				New: changedFile,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.FileUpdated(); got != tt.want {
				t.Errorf("Scanned.FileUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanned_ContentsChanged(t *testing.T) {
	const (
		checksum = "checksum"
		oshash   = "oshash"
		changed  = "changed"
	)

	tests := []struct {
		name string
		s    Scanned
		want bool
	}{
		{
			"unchanged",
			Scanned{
				Old: &models.File{
					Checksum: checksum,
					OSHash:   oshash,
				},
				New: &models.File{
					Checksum: checksum,
					OSHash:   oshash,
				},
			},
			false,
		},
		{
			"checksum changed",
			Scanned{
				Old: &models.File{
					Checksum: checksum,
					OSHash:   oshash,
				},
				New: &models.File{
					Checksum: changed,
					OSHash:   oshash,
				},
			},
			true,
		},
		{
			"oshash changed",
			Scanned{
				Old: &models.File{
					Checksum: checksum,
					OSHash:   oshash,
				},
				New: &models.File{
					Checksum: checksum,
					OSHash:   changed,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ContentsChanged(); got != tt.want {
				t.Errorf("Scanned.ContentsChanged() = %v, want %v", got, tt.want)
			}
		})
	}
}

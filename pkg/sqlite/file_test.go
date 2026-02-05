//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
)

func getFilePath(folderIdx int, basename string) string {
	return filepath.Join(folderPaths[folderIdx], basename)
}

func makeZipFileWithID(index int) models.File {
	f := makeFile(index)

	return &models.BaseFile{
		ID:       fileIDs[index],
		Basename: f.Base().Basename,
		Path:     getFilePath(fileFolders[index], getFileBaseName(index)),
	}
}

func Test_fileFileStore_Create(t *testing.T) {
	var (
		basename               = "basename"
		fingerprintType        = "MD5"
		fingerprintValue       = "checksum"
		fileModTime            = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt              = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt              = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		size             int64 = 1234

		duration         = 1.234
		width            = 640
		height           = 480
		framerate        = 2.345
		bitrate    int64 = 234
		videoCodec       = "videoCodec"
		audioCodec       = "audioCodec"
		format           = "format"
	)

	tests := []struct {
		name      string
		newObject models.File
		wantErr   bool
	}{
		{
			"full",
			&models.BaseFile{
				DirEntry: models.DirEntry{
					ZipFileID: &fileIDs[fileIdxZip],
					ZipFile:   makeZipFileWithID(fileIdxZip),
					ModTime:   fileModTime,
				},
				Path:           getFilePath(folderIdxWithFiles, basename),
				ParentFolderID: folderIDs[folderIdxWithFiles],
				Basename:       basename,
				Size:           size,
				Fingerprints: []models.Fingerprint{
					{
						Type:        fingerprintType,
						Fingerprint: fingerprintValue,
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"video file",
			&models.VideoFile{
				BaseFile: &models.BaseFile{
					DirEntry: models.DirEntry{
						ZipFileID: &fileIDs[fileIdxZip],
						ZipFile:   makeZipFileWithID(fileIdxZip),
						ModTime:   fileModTime,
					},
					Path:           getFilePath(folderIdxWithFiles, basename),
					ParentFolderID: folderIDs[folderIdxWithFiles],
					Basename:       basename,
					Size:           size,
					Fingerprints: []models.Fingerprint{
						{
							Type:        fingerprintType,
							Fingerprint: fingerprintValue,
						},
					},
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Duration:   duration,
				VideoCodec: videoCodec,
				AudioCodec: audioCodec,
				Format:     format,
				Width:      width,
				Height:     height,
				FrameRate:  framerate,
				BitRate:    bitrate,
			},
			false,
		},
		{
			"image file",
			&models.ImageFile{
				BaseFile: &models.BaseFile{
					DirEntry: models.DirEntry{
						ZipFileID: &fileIDs[fileIdxZip],
						ZipFile:   makeZipFileWithID(fileIdxZip),
						ModTime:   fileModTime,
					},
					Path:           getFilePath(folderIdxWithFiles, basename),
					ParentFolderID: folderIDs[folderIdxWithFiles],
					Basename:       basename,
					Size:           size,
					Fingerprints: []models.Fingerprint{
						{
							Type:        fingerprintType,
							Fingerprint: fingerprintValue,
						},
					},
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Format: format,
				Width:  width,
				Height: height,
			},
			false,
		},
		{
			"duplicate path",
			&models.BaseFile{
				DirEntry: models.DirEntry{
					ModTime: fileModTime,
				},
				Path:           getFilePath(folderIdxWithFiles, getFileBaseName(fileIdxZip)),
				ParentFolderID: folderIDs[folderIdxWithFiles],
				Basename:       getFileBaseName(fileIdxZip),
				Size:           size,
				Fingerprints: []models.Fingerprint{
					{
						Type:        fingerprintType,
						Fingerprint: fingerprintValue,
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"empty basename",
			&models.BaseFile{
				ParentFolderID: folderIDs[folderIdxWithFiles],
			},
			true,
		},
		{
			"missing folder id",
			&models.BaseFile{
				Basename: basename,
			},
			true,
		},
		{
			"invalid folder id",
			&models.BaseFile{
				DirEntry:       models.DirEntry{},
				ParentFolderID: invalidFolderID,
				Basename:       basename,
			},
			true,
		},
		{
			"invalid zip file id",
			&models.BaseFile{
				DirEntry: models.DirEntry{
					ZipFileID: &invalidFileID,
				},
				Basename: basename,
			},
			true,
		},
	}

	qb := db.File

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			s := tt.newObject
			if err := qb.Create(ctx, s); (err != nil) != tt.wantErr {
				t.Errorf("fileStore.Create() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				assert.Zero(s.Base().ID)
				return
			}

			assert.NotZero(s.Base().ID)

			var copy models.File
			switch t := s.(type) {
			case *models.BaseFile:
				v := *t
				copy = &v
			case *models.VideoFile:
				v := *t
				copy = &v
			case *models.ImageFile:
				v := *t
				copy = &v
			}

			copy.Base().ID = s.Base().ID

			assert.Equal(copy, s)

			// ensure can find the scene
			found, err := qb.Find(ctx, s.Base().ID)
			if err != nil {
				t.Errorf("fileStore.Find() error = %v", err)
			}

			if !assert.Len(found, 1) {
				return
			}

			assert.Equal(copy, found[0])

			return
		})
	}
}

func Test_fileStore_Update(t *testing.T) {
	var (
		basename               = "basename"
		fingerprintType        = "MD5"
		fingerprintValue       = "checksum"
		fileModTime            = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		createdAt              = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt              = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
		size             int64 = 1234

		duration         = 1.234
		width            = 640
		height           = 480
		framerate        = 2.345
		bitrate    int64 = 234
		videoCodec       = "videoCodec"
		audioCodec       = "audioCodec"
		format           = "format"
	)

	tests := []struct {
		name          string
		updatedObject models.File
		wantErr       bool
	}{
		{
			"full",
			&models.BaseFile{
				ID: fileIDs[fileIdxInZip],
				DirEntry: models.DirEntry{
					ZipFileID: &fileIDs[fileIdxZip],
					ZipFile:   makeZipFileWithID(fileIdxZip),
					ModTime:   fileModTime,
				},
				Path:           getFilePath(folderIdxWithFiles, basename),
				ParentFolderID: folderIDs[folderIdxWithFiles],
				Basename:       basename,
				Size:           size,
				Fingerprints: []models.Fingerprint{
					{
						Type:        fingerprintType,
						Fingerprint: fingerprintValue,
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			false,
		},
		{
			"video file",
			&models.VideoFile{
				BaseFile: &models.BaseFile{
					ID: fileIDs[fileIdxStartVideoFiles],
					DirEntry: models.DirEntry{
						ZipFileID: &fileIDs[fileIdxZip],
						ZipFile:   makeZipFileWithID(fileIdxZip),
						ModTime:   fileModTime,
					},
					Path:           getFilePath(folderIdxWithFiles, basename),
					ParentFolderID: folderIDs[folderIdxWithFiles],
					Basename:       basename,
					Size:           size,
					Fingerprints: []models.Fingerprint{
						{
							Type:        fingerprintType,
							Fingerprint: fingerprintValue,
						},
					},
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Duration:   duration,
				VideoCodec: videoCodec,
				AudioCodec: audioCodec,
				Format:     format,
				Width:      width,
				Height:     height,
				FrameRate:  framerate,
				BitRate:    bitrate,
			},
			false,
		},
		{
			"image file",
			&models.ImageFile{
				BaseFile: &models.BaseFile{
					ID: fileIDs[fileIdxStartImageFiles],
					DirEntry: models.DirEntry{
						ZipFileID: &fileIDs[fileIdxZip],
						ZipFile:   makeZipFileWithID(fileIdxZip),
						ModTime:   fileModTime,
					},
					Path:           getFilePath(folderIdxWithFiles, basename),
					ParentFolderID: folderIDs[folderIdxWithFiles],
					Basename:       basename,
					Size:           size,
					Fingerprints: []models.Fingerprint{
						{
							Type:        fingerprintType,
							Fingerprint: fingerprintValue,
						},
					},
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				Format: format,
				Width:  width,
				Height: height,
			},
			false,
		},
		{
			"duplicate path",
			&models.BaseFile{
				ID: fileIDs[fileIdxInZip],
				DirEntry: models.DirEntry{
					ModTime: fileModTime,
				},
				Path:           getFilePath(folderIdxWithFiles, getFileBaseName(fileIdxZip)),
				ParentFolderID: folderIDs[folderIdxWithFiles],
				Basename:       getFileBaseName(fileIdxZip),
				Size:           size,
				Fingerprints: []models.Fingerprint{
					{
						Type:        fingerprintType,
						Fingerprint: fingerprintValue,
					},
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			true,
		},
		{
			"clear zip",
			&models.BaseFile{
				ID:             fileIDs[fileIdxInZip],
				Path:           getFilePath(folderIdxWithFiles, getFileBaseName(fileIdxZip)+".renamed"),
				Basename:       getFileBaseName(fileIdxZip) + ".renamed",
				ParentFolderID: folderIDs[folderIdxWithFiles],
			},
			false,
		},
		{
			"clear folder",
			&models.BaseFile{
				ID:   fileIDs[fileIdxZip],
				Path: basename,
			},
			true,
		},
		{
			"invalid parent folder id",
			&models.BaseFile{
				ID:             fileIDs[fileIdxZip],
				Path:           basename,
				ParentFolderID: invalidFolderID,
			},
			true,
		},
		{
			"invalid zip file id",
			&models.BaseFile{
				ID:   fileIDs[fileIdxZip],
				Path: basename,
				DirEntry: models.DirEntry{
					ZipFileID: &invalidFileID,
				},
				ParentFolderID: folderIDs[folderIdxWithFiles],
			},
			true,
		},
	}

	qb := db.File
	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)

			copy := tt.updatedObject

			if err := qb.Update(ctx, tt.updatedObject); (err != nil) != tt.wantErr {
				t.Errorf("FileStore.Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			s, err := qb.Find(ctx, tt.updatedObject.Base().ID)
			if err != nil {
				t.Errorf("FileStore.Find() error = %v", err)
			}

			if !assert.Len(s, 1) {
				return
			}

			assert.Equal(copy, s[0])

			return
		})
	}
}

func makeFileWithID(index int) models.File {
	ret := makeFile(index)
	ret.Base().Path = getFilePath(fileFolders[index], getFileBaseName(index))
	ret.Base().ID = fileIDs[index]

	return ret
}

func Test_fileStore_Find(t *testing.T) {
	tests := []struct {
		name    string
		id      models.FileID
		want    models.File
		wantErr bool
	}{
		{
			"valid",
			fileIDs[fileIdxZip],
			makeFileWithID(fileIdxZip),
			false,
		},
		{
			"invalid",
			models.FileID(invalidID),
			nil,
			true,
		},
		{
			"video file",
			fileIDs[fileIdxStartVideoFiles],
			makeFileWithID(fileIdxStartVideoFiles),
			false,
		},
		{
			"image file",
			fileIDs[fileIdxStartImageFiles],
			makeFileWithID(fileIdxStartImageFiles),
			false,
		},
	}

	qb := db.File

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.Find(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("fileStore.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil {
				assert.Len(got, 0)
				return
			}

			if !assert.Len(got, 1) {
				return
			}

			assert.Equal(tt.want, got[0])
		})
	}
}

func Test_FileStore_FindByPath(t *testing.T) {
	getPath := func(index int) string {
		folderIdx, found := fileFolders[index]
		if !found {
			folderIdx = folderIdxWithFiles
		}

		return getFilePath(folderIdx, getFileBaseName(index))
	}

	tests := []struct {
		name    string
		path    string
		want    models.File
		wantErr bool
	}{
		{
			"valid",
			getPath(fileIdxZip),
			makeFileWithID(fileIdxZip),
			false,
		},
		{
			"invalid",
			"invalid path",
			nil,
			false,
		},
	}

	qb := db.File

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByPath(ctx, tt.path, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.FindByPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func TestFileStore_FindByFingerprint(t *testing.T) {
	tests := []struct {
		name    string
		fp      models.Fingerprint
		want    []models.File
		wantErr bool
	}{
		{
			"by MD5",
			models.Fingerprint{
				Type:        "MD5",
				Fingerprint: getPrefixedStringValue("file", fileIdxZip, "md5"),
			},
			[]models.File{makeFileWithID(fileIdxZip)},
			false,
		},
		{
			"by OSHASH",
			models.Fingerprint{
				Type:        "OSHASH",
				Fingerprint: getPrefixedStringValue("file", fileIdxZip, "oshash"),
			},
			[]models.File{makeFileWithID(fileIdxZip)},
			false,
		},
		{
			"non-existing",
			models.Fingerprint{
				Type:        "OSHASH",
				Fingerprint: "foo",
			},
			nil,
			false,
		},
	}

	qb := db.File

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.FindByFingerprint(ctx, tt.fp)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.FindByFingerprint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

func TestFileStore_IsPrimary(t *testing.T) {
	tests := []struct {
		name   string
		fileID models.FileID
		want   bool
	}{
		{
			"scene file",
			sceneFileIDs[sceneIdx1WithPerformer],
			true,
		},
		{
			"image file",
			imageFileIDs[imageIdx1WithGallery],
			true,
		},
		{
			"gallery file",
			galleryFileIDs[galleryIdx1WithImage],
			true,
		},
		{
			"orphan file",
			fileIDs[fileIdxZip],
			false,
		},
		{
			"invalid file",
			invalidFileID,
			false,
		},
	}

	qb := db.File

	for _, tt := range tests {
		runWithRollbackTxn(t, tt.name, func(t *testing.T, ctx context.Context) {
			assert := assert.New(t)
			got, err := qb.IsPrimary(ctx, tt.fileID)
			if err != nil {
				t.Errorf("FileStore.IsPrimary() error = %v", err)
				return
			}

			assert.Equal(tt.want, got)
		})
	}
}

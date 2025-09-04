package models

import (
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"time"
)

// FolderID represents an ID of a folder.
type FolderID int32

// String converts the ID to a string.
func (i FolderID) String() string {
	return strconv.Itoa(int(i))
}

func (i *FolderID) UnmarshalGQL(v interface{}) (err error) {
	switch v := v.(type) {
	case string:
		var id int
		id, err = strconv.Atoi(v)
		*i = FolderID(id)
		return err
	case int:
		*i = FolderID(v)
		return nil
	default:
		return fmt.Errorf("%T is not an int", v)
	}
}

func (i FolderID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(i.String()))
}

func FolderIDsFromInts(ids []int) []FolderID {
	ret := make([]FolderID, len(ids))
	for i, id := range ids {
		ret[i] = FolderID(id)
	}
	return ret
}

// Folder represents a folder in the file system.
type Folder struct {
	ID FolderID `json:"id"`
	DirEntry
	Path           string    `json:"path"`
	ParentFolderID *FolderID `json:"parent_folder_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Folder) Info(fs FS) (fs.FileInfo, error) {
	return f.info(fs, f.Path)
}

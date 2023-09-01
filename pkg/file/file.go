package file

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// Repository provides access to storage methods for files and folders.
type Repository struct {
	txn.Manager
	txn.DatabaseProvider

	FileStore   models.FileReaderWriter
	FolderStore models.FolderReaderWriter
}

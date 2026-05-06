package manager

import (
	"github.com/stashapp/stash/internal/manager/config"
)

type SystemStatus struct {
	DatabaseSchema *int             `json:"databaseSchema"`
	DatabasePath   *string          `json:"databasePath"`
	ConfigPath     *string          `json:"configPath"`
	AppSchema      int              `json:"appSchema"`
	Status         SystemStatusEnum `json:"status"`
	Os             string           `json:"os"`
	WorkingDir     string           `json:"working_dir"`
	HomeDir        string           `json:"home_dir"`
	FfmpegPath     *string          `json:"ffmpegPath"`
	FfprobePath    *string          `json:"ffprobePath"`
}

type SetupInput struct {
	// Empty to indicate $HOME/.stash/config.yml default
	ConfigLocation string                     `json:"configLocation"`
	Stashes        []*config.StashConfigInput `json:"stashes"`
	SFWContentMode bool                       `json:"sfwContentMode"`
	// Empty to indicate default
	DatabaseFile string `json:"databaseFile"`
	// Empty to indicate default
	GeneratedLocation string `json:"generatedLocation"`
	// Empty to indicate default
	CacheLocation string `json:"cacheLocation"`

	StoreBlobsInDatabase bool `json:"storeBlobsInDatabase"`
	// Empty to indicate default
	BlobsLocation string `json:"blobsLocation"`
}

type MigrateInput struct {
	BackupPath string `json:"backupPath"`
}

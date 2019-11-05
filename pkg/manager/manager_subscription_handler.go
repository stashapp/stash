package manager

import (
	"github.com/stashapp/stash/pkg/models"
)

func (s *singleton) GetMetadataUpdateStatus() models.MetadataUpdateStatus {
	ret := models.MetadataUpdateStatus{
		Progress: instance.Status.Progress,
		Status:   instance.Status.Status.String(),
		Message:  "",
	}
	return ret
}

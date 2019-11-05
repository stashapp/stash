package manager

import (
	"github.com/stashapp/stash/pkg/models"
)

func (s *singleton) GetMetadataUpdateStatus() models.MetadataUpdateStatus {
	ret := models.MetadataUpdateStatus{
		Progress: -1,
		Status:   instance.Status.String(),
		Message:  "",
	}
	return ret
}

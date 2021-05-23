package api

import (
	"context"
	"time"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
)

func (r *subscriptionResolver) MetadataUpdate(ctx context.Context) (<-chan *models.MetadataUpdateStatus, error) {
	msg := make(chan *models.MetadataUpdateStatus, 1)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		lastStatus := manager.TaskStatus{}
		for {
			select {
			case <-ticker.C:
				thisStatus := manager.GetInstance().Status
				if thisStatus != lastStatus {
					ret := models.MetadataUpdateStatus{
						Progress: thisStatus.Progress,
						Status:   thisStatus.Status.String(),
						Message:  "",
					}
					msg <- &ret
				}
				lastStatus = thisStatus
			case <-ctx.Done():
				ticker.Stop()
				close(msg)
				return
			}
		}
	}()

	return msg, nil
}

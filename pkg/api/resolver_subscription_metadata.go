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
		lastStatus := models.MetadataUpdateStatus{}
		for {
			select {
			case _ = <-ticker.C:
				thisStatus := manager.GetInstance().GetMetadataUpdateStatus()
				if thisStatus != lastStatus {
					msg <- &thisStatus
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

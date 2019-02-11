package api

import (
	"context"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/manager"
	"time"
)

func (r *subscriptionResolver) MetadataUpdate(ctx context.Context) (<-chan string, error) {
	msg := make(chan string, 1)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				logger.Trace("metadata subscription tick at %s", t)
				manager.GetInstance().HandleMetadataUpdateSubscriptionTick(msg)
			case <-ctx.Done():
				ticker.Stop()
				close(msg)
				return
			}
		}
	}()

	return msg, nil
}
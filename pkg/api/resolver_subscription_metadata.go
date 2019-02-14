package api

import (
	"context"
	"github.com/stashapp/stash/pkg/manager"
	"time"
)

func (r *subscriptionResolver) MetadataUpdate(ctx context.Context) (<-chan string, error) {
	msg := make(chan string, 1)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case _ = <-ticker.C:
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

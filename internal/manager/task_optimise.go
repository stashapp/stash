package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
)

type Optimiser interface {
	Analyze(ctx context.Context) error
	Vacuum(ctx context.Context) error
}

type OptimiseDatabaseJob struct {
	Optimiser Optimiser
}

func (j *OptimiseDatabaseJob) Execute(ctx context.Context, progress *job.Progress) error {
	logger.Info("Optimising database")
	progress.SetTotal(2)

	start := time.Now()

	var err error

	progress.ExecuteTask("Analyzing database", func() {
		err = j.Optimiser.Analyze(ctx)
		progress.Increment()
	})
	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error analyzing database: %w", err)
	}

	progress.ExecuteTask("Vacuuming database", func() {
		err = j.Optimiser.Vacuum(ctx)
		progress.Increment()
	})
	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error vacuuming database: %w", err)
	}

	elapsed := time.Since(start)
	logger.Infof("Finished optimising database after %s", elapsed)
	return nil
}

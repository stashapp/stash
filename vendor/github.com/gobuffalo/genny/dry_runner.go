package genny

import (
	"context"

	"github.com/gobuffalo/logger"
)

// DryRunner will NOT execute commands and write files
// it is NOT destructive
func DryRunner(ctx context.Context) *Runner {
	r := NewRunner(ctx)
	r.Logger = logger.New(logger.DebugLevel)
	return r
}

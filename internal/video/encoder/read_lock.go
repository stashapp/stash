package encoder

import "github.com/stashapp/stash/pkg/fsutil"

// TODO - this should be a dependency passed through, rather than a global variable.

var readLockManager *fsutil.ReadLockManager = fsutil.NewReadLockManager()

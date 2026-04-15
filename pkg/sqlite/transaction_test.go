//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
)

// this test is left commented out as it is not deterministic.
// func TestConcurrentExclusiveTxn(t *testing.T) {
// 	const (
// 		workers    = 8
// 		loops      = 100
// 		innerLoops = 10
// 		sleepTime  = 2 * time.Millisecond
// 	)
// 	ctx := context.Background()

// 	var wg sync.WaitGroup
// 	for k := 0; k < workers; k++ {
// 		wg.Add(1)
// 		go func(wk int) {
// 			for l := 0; l < loops; l++ {
// 				// change this to WithReadTxn to see locked database error
// 				if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
// 					for ll := 0; ll < innerLoops; ll++ {
// 						scene := &models.Scene{
// 							Title: "test",
// 						}

// 						if err := db.Scene.Create(ctx, scene, nil); err != nil {
// 							return err
// 						}

// 						if err := db.Scene.Destroy(ctx, scene.ID); err != nil {
// 							return err
// 						}
// 					}
// 					time.Sleep(sleepTime)

// 					return nil
// 				}); err != nil {
// 					t.Errorf("worker %d loop %d: %v", wk, l, err)
// 				}
// 			}

// 			wg.Done()
// 		}(k)
// 	}

// 	wg.Wait()
// }

func signalOtherThread(c chan struct{}) error {
	select {
	case c <- struct{}{}:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("timed out signalling other thread")
	}
}

func waitForOtherThread(c chan struct{}) error {
	select {
	case <-c:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("timed out waiting for other thread")
	}
}

// this test is left commented as it's no longer possible to write to the database
// with a read-only transaction.

// func TestConcurrentReadTxn(t *testing.T) {
// 	var wg sync.WaitGroup
// 	ctx := context.Background()
// 	c := make(chan struct{})

// 	// first thread
// 	wg.Add(2)
// 	go func() {
// 		defer wg.Done()
// 		if err := txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
// 			scene := &models.Scene{
// 				Title: "test",
// 			}

// 			if err := db.Scene.Create(ctx, scene, nil); err != nil {
// 				return err
// 			}

// 			// wait for other thread to start
// 			if err := signalOtherThread(c); err != nil {
// 				return err
// 			}
// 			if err := waitForOtherThread(c); err != nil {
// 				return err
// 			}

// 			if err := db.Scene.Destroy(ctx, scene.ID); err != nil {
// 				return err
// 			}

// 			return nil
// 		}); err != nil {
// 			t.Errorf("unexpected error in first thread: %v", err)
// 		}
// 	}()

// 	// second thread
// 	go func() {
// 		defer wg.Done()
// 		_ = txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
// 			// wait for first thread
// 			if err := waitForOtherThread(c); err != nil {
// 				t.Errorf(err.Error())
// 				return err
// 			}

// 			defer func() {
// 				if err := signalOtherThread(c); err != nil {
// 					t.Errorf(err.Error())
// 				}
// 			}()

// 			scene := &models.Scene{
// 				Title: "test",
// 			}

// 			// expect error when we try to do this, as the other thread has already
// 			// modified this table
// 			// this takes time to fail, so we need to wait for it
// 			if err := db.Scene.Create(ctx, scene, nil); err != nil {
// 				if !db.IsLocked(err) {
// 					t.Errorf("unexpected error: %v", err)
// 				}
// 				return err
// 			} else {
// 				t.Errorf("expected locked error in second thread")
// 			}

// 			return nil
// 		})
// 	}()

// 	wg.Wait()
// }

func TestConcurrentExclusiveAndReadTxn(t *testing.T) {
	var wg sync.WaitGroup
	ctx := context.Background()
	c := make(chan struct{})

	// first thread
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
			scene := &models.Scene{
				Title: "test",
			}

			if err := db.Scene.Create(ctx, scene, nil); err != nil {
				return err
			}

			// wait for other thread to start
			if err := signalOtherThread(c); err != nil {
				return err
			}
			if err := waitForOtherThread(c); err != nil {
				return err
			}

			if err := db.Scene.Destroy(ctx, scene.ID); err != nil {
				return err
			}

			return nil
		}); err != nil {
			t.Errorf("unexpected error in first thread: %v", err)
		}
	}()

	// second thread
	go func() {
		defer wg.Done()
		_ = txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
			// wait for first thread
			if err := waitForOtherThread(c); err != nil {
				t.Error(err.Error())
				return err
			}

			defer func() {
				if err := signalOtherThread(c); err != nil {
					t.Error(err.Error())
				}
			}()

			if _, err := db.Scene.Find(ctx, sceneIDs[sceneIdx1WithPerformer]); err != nil {
				t.Errorf("unexpected error: %v", err)
				return err
			}

			return nil
		})
	}()

	wg.Wait()
}

// this test is left commented out as it is not deterministic.
// func TestConcurrentExclusiveAndReadTxns(t *testing.T) {
// 	const (
// 		writeWorkers = 4
// 		readWorkers  = 4
// 		loops        = 200
// 		innerLoops   = 10
// 		sleepTime  = 1 * time.Millisecond
// 	)
// 	ctx := context.Background()

// 	var wg sync.WaitGroup
// 	for k := 0; k < writeWorkers; k++ {
// 		wg.Add(1)
// 		go func(wk int) {
// 			for l := 0; l < loops; l++ {
// 				if err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
// 					for ll := 0; ll < innerLoops; ll++ {
// 						scene := &models.Scene{
// 							Title: "test",
// 						}

// 						if err := db.Scene.Create(ctx, scene, nil); err != nil {
// 							return err
// 						}

// 						if err := db.Scene.Destroy(ctx, scene.ID); err != nil {
// 							return err
// 						}
// 					}
// 					time.Sleep(sleepTime)

// 					return nil
// 				}); err != nil {
// 					t.Errorf("write worker %d loop %d: %v", wk, l, err)
// 				}
// 			}

// 			wg.Done()
// 		}(k)
// 	}

// 	for k := 0; k < readWorkers; k++ {
// 		wg.Add(1)
// 		go func(wk int) {
// 			for l := 0; l < loops; l++ {
// 				if err := txn.WithReadTxn(ctx, db, func(ctx context.Context) error {
// 					for ll := 0; ll < innerLoops; ll++ {
// 						if _, err := db.Scene.Find(ctx, sceneIDs[ll%totalScenes]); err != nil {
// 							return err
// 						}
// 					}
// 					time.Sleep(sleepTime)

// 					return nil
// 				}); err != nil {
// 					t.Errorf("read worker %d loop %d: %v", wk, l, err)
// 				}
// 			}

// 			wg.Done()
// 		}(k)
// 	}

// 	wg.Wait()
// }

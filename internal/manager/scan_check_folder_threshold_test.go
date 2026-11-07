package manager

import (
	"testing"

	"github.com/stashapp/stash/pkg/file"
)

func newTestJob(threshold float64, skipUnchanged bool, rescan bool) *ScanJob {
	j := &ScanJob{
		scanner:              &file.Scanner{Rescan: rescan},
		checkFolderThreshold: threshold,
		skipUnchangedFolders: skipUnchanged,
	}
	return j
}

func simulateChecks(j *ScanJob, hits, misses int) {
	for i := 0; i < hits; i++ {
		j.dirCheckAttempts.Add(1)
		j.dirCheckHits.Add(1)
	}
	for i := 0; i < misses; i++ {
		j.dirCheckAttempts.Add(1)
	}
}

// TestShouldCheckFolder_WarmupWindow verifies that the first dirCheckWarmup
// calls always return true regardless of hit rate.
func TestShouldCheckFolder_WarmupWindow(t *testing.T) {
	j := newTestJob(0.30, true, false)

	// Simulate all misses but still within warmup window.
	for i := 0; i < dirCheckWarmup-1; i++ {
		j.dirCheckAttempts.Add(1) // all misses
		if !j.shouldCheckFolder() {
			t.Fatalf("shouldCheckFolder returned false at attempt %d (within warmup window)", i+1)
		}
	}
}

// TestShouldCheckFolder_AllHits verifies that a 100% hit rate keeps
// shouldCheckFolder enabled well past the warmup window.
func TestShouldCheckFolder_AllHits(t *testing.T) {
	j := newTestJob(0.30, true, false)
	simulateChecks(j, 200, 0)

	for i := 0; i < 10; i++ {
		if !j.shouldCheckFolder() {
			t.Fatal("shouldCheckFolder returned false with 100% hit rate")
		}
	}
}

// TestShouldCheckFolder_AllMisses verifies that a 0% hit rate disables
// shouldCheckFolder after the warmup window.
func TestShouldCheckFolder_AllMisses(t *testing.T) {
	j := newTestJob(0.30, true, false)
	simulateChecks(j, 0, dirCheckWarmup+10)

	if j.shouldCheckFolder() {
		t.Fatal("shouldCheckFolder returned true with 0% hit rate after warmup")
	}
}

// TestShouldCheckFolder_Transition verifies the transition: starts warm (high
// hit rate, enabled), then accumulates misses until the rate drops below
// threshold and the gate disables.
func TestShouldCheckFolder_Transition(t *testing.T) {
	j := newTestJob(0.30, true, false)

	// Start warm: fill warmup window with hits.
	simulateChecks(j, dirCheckWarmup, 0)
	if !j.shouldCheckFolder() {
		t.Fatal("shouldCheckFolder returned false with 100% hit rate after warmup")
	}

	// Add enough misses to push hit rate below threshold (0.30).
	// After 50 hits + N misses, rate = 50 / (50+N). Want rate < 0.30:
	// 50 / (50+N) < 0.30 → N > 50/0.30 - 50 ≈ 117.
	simulateChecks(j, 0, 120)
	if j.shouldCheckFolder() {
		total := j.dirCheckAttempts.Load()
		hits := j.dirCheckHits.Load()
		t.Fatalf("shouldCheckFolder returned true with hit rate %.2f (below threshold 0.30), total=%d hits=%d",
			float64(hits)/float64(total), total, hits)
	}
}

// TestShouldCheckFolder_RescanBypass verifies that Rescan=true always
// returns false, regardless of hit rate or warmup state.
func TestShouldCheckFolder_RescanBypass(t *testing.T) {
	j := newTestJob(0.30, true, true) // Rescan=true
	simulateChecks(j, 200, 0)         // 100% hits, past warmup

	if j.shouldCheckFolder() {
		t.Fatal("shouldCheckFolder returned true when Rescan=true")
	}
}

// TestShouldCheckFolder_NetworkFSBypass verifies that skipUnchangedFolders=false
// (network filesystem) always returns false.
func TestShouldCheckFolder_NetworkFSBypass(t *testing.T) {
	j := newTestJob(0.30, false, false) // skipUnchangedFolders=false
	simulateChecks(j, 200, 0)

	if j.shouldCheckFolder() {
		t.Fatal("shouldCheckFolder returned true when skipUnchangedFolders=false")
	}
}

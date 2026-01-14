package dlna

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/txn"
)

const (
	// DefaultSessionTimeout is the time after which a session is considered complete
	// if no new requests are received.
	// This is set high (5 minutes) because DLNA clients buffer aggressively and may not
	// send any HTTP requests for extended periods while the user is still watching.
	DefaultSessionTimeout = 5 * time.Minute

	// monitorInterval is how often we check for expired sessions.
	monitorInterval = 10 * time.Second
)

// ActivityConfig provides configuration options for DLNA activity tracking.
type ActivityConfig interface {
	// GetDLNAActivityTrackingEnabled returns true if activity tracking should be enabled.
	// If not implemented, defaults to true.
	GetDLNAActivityTrackingEnabled() bool

	// GetMinimumPlayPercent returns the minimum percentage of a video that must be
	// watched before incrementing the play count. Uses UI setting if available.
	GetMinimumPlayPercent() int
}

// SceneActivityWriter provides methods for saving scene activity.
type SceneActivityWriter interface {
	SaveActivity(ctx context.Context, sceneID int, resumeTime *float64, playDuration *float64) (bool, error)
	AddViews(ctx context.Context, sceneID int, dates []time.Time) ([]time.Time, error)
}

// streamSession represents an active DLNA streaming session.
type streamSession struct {
	SceneID        int
	ClientIP       string
	StartTime      time.Time
	LastActivity   time.Time
	VideoDuration  float64
	PlayCountAdded bool
}

// sessionKey generates a unique key for a session based on client IP and scene ID.
func sessionKey(clientIP string, sceneID int) string {
	return fmt.Sprintf("%s:%d", clientIP, sceneID)
}

// percentWatched calculates the estimated percentage of video watched.
// Uses a time-based approach since DLNA clients buffer aggressively and byte
// positions don't correlate with actual playback position.
//
// The key insight: you cannot have watched more of the video than time has elapsed.
// If the video is 30 minutes and only 1 minute has passed, maximum watched is ~3.3%.
func (s *streamSession) percentWatched() float64 {
	if s.VideoDuration <= 0 {
		return 0
	}

	// Calculate elapsed time from session start to last activity
	elapsed := s.LastActivity.Sub(s.StartTime).Seconds()
	if elapsed <= 0 {
		return 0
	}

	// Maximum possible percent is based on elapsed time
	// You can't watch more of the video than time has passed
	timeBasedPercent := (elapsed / s.VideoDuration) * 100

	// Cap at 100%
	if timeBasedPercent > 100 {
		return 100
	}

	return timeBasedPercent
}

// estimatedResumeTime calculates the estimated resume time based on elapsed time.
// Since DLNA clients buffer aggressively, byte positions don't correlate with playback.
// Instead, we estimate based on how long the session has been active.
// Returns the time in seconds, or 0 if the video is nearly complete (>=98%).
func (s *streamSession) estimatedResumeTime() float64 {
	if s.VideoDuration <= 0 {
		return 0
	}

	// Calculate elapsed time from session start
	elapsed := s.LastActivity.Sub(s.StartTime).Seconds()
	if elapsed <= 0 {
		return 0
	}

	// If elapsed time exceeds 98% of video duration, reset resume time (matches frontend behavior)
	if elapsed >= s.VideoDuration*0.98 {
		return 0
	}

	// Resume time is approximately where the user was watching
	// Capped by video duration
	if elapsed > s.VideoDuration {
		elapsed = s.VideoDuration
	}

	return elapsed
}

// ActivityTracker tracks DLNA streaming activity and saves it to the database.
type ActivityTracker struct {
	txnManager     txn.Manager
	sceneWriter    SceneActivityWriter
	config         ActivityConfig
	sessionTimeout time.Duration

	sessions map[string]*streamSession
	mutex    sync.RWMutex

	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// NewActivityTracker creates a new ActivityTracker.
func NewActivityTracker(
	txnManager txn.Manager,
	sceneWriter SceneActivityWriter,
	config ActivityConfig,
) *ActivityTracker {
	ctx, cancel := context.WithCancel(context.Background())

	tracker := &ActivityTracker{
		txnManager:     txnManager,
		sceneWriter:    sceneWriter,
		config:         config,
		sessionTimeout: DefaultSessionTimeout,
		sessions:       make(map[string]*streamSession),
		ctx:            ctx,
		cancelFunc:     cancel,
	}

	// Start the session monitor goroutine
	tracker.wg.Add(1)
	go tracker.monitorSessions()

	return tracker
}

// Stop stops the activity tracker and processes any remaining sessions.
func (t *ActivityTracker) Stop() {
	t.cancelFunc()
	t.wg.Wait()

	// Process any remaining sessions
	t.mutex.Lock()
	sessions := make([]*streamSession, 0, len(t.sessions))
	for _, session := range t.sessions {
		sessions = append(sessions, session)
	}
	t.sessions = make(map[string]*streamSession)
	t.mutex.Unlock()

	for _, session := range sessions {
		t.processCompletedSession(session)
	}
}

// RecordRequest records a streaming request for activity tracking.
// Each request updates the session's LastActivity time, which is used for
// time-based tracking of watch progress.
func (t *ActivityTracker) RecordRequest(sceneID int, clientIP string, videoDuration float64) {
	if !t.isEnabled() {
		return
	}

	key := sessionKey(clientIP, sceneID)
	now := time.Now()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	session, exists := t.sessions[key]
	if !exists {
		session = &streamSession{
			SceneID:       sceneID,
			ClientIP:      clientIP,
			StartTime:     now,
			VideoDuration: videoDuration,
		}
		t.sessions[key] = session
		logger.Debugf("[DLNA Activity] New session started: scene=%d, client=%s", sceneID, clientIP)
	}

	session.LastActivity = now
}

// monitorSessions periodically checks for expired sessions and processes them.
func (t *ActivityTracker) monitorSessions() {
	defer t.wg.Done()

	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			return
		case <-ticker.C:
			t.processExpiredSessions()
		}
	}
}

// processExpiredSessions finds and processes sessions that have timed out.
func (t *ActivityTracker) processExpiredSessions() {
	now := time.Now()
	var expiredSessions []*streamSession

	t.mutex.Lock()
	for key, session := range t.sessions {
		timeSinceStart := now.Sub(session.StartTime)
		timeSinceActivity := now.Sub(session.LastActivity)

		// Must have no HTTP activity for the full timeout period
		if timeSinceActivity <= t.sessionTimeout {
			continue
		}

		// DLNA clients buffer aggressively - they fetch most/all of the video quickly,
		// then play from cache with NO further HTTP requests.
		//
		// Two scenarios:
		// 1. User watched the whole video: timeSinceStart >= videoDuration
		//    -> Set LastActivity to when timeout began (they finished watching)
		// 2. User stopped early: timeSinceStart < videoDuration
		//    -> Keep LastActivity as-is (best estimate of when they stopped)

		videoDuration := time.Duration(session.VideoDuration) * time.Second
		if timeSinceStart >= videoDuration && videoDuration > 0 {
			// User likely watched the whole video, then it timed out
			// Estimate they watched until the timeout period started
			session.LastActivity = now.Add(-t.sessionTimeout)
		}
		// else: User stopped early - LastActivity is already our best estimate

		expiredSessions = append(expiredSessions, session)
		delete(t.sessions, key)
	}
	t.mutex.Unlock()

	for _, session := range expiredSessions {
		t.processCompletedSession(session)
	}
}

// processCompletedSession saves activity data for a completed streaming session.
func (t *ActivityTracker) processCompletedSession(session *streamSession) {
	percentWatched := session.percentWatched()
	resumeTime := session.estimatedResumeTime()

	logger.Debugf("[DLNA Activity] Session completed: scene=%d, client=%s, videoDuration=%.1fs, percent=%.1f%%, resume=%.1fs",
		session.SceneID, session.ClientIP, session.VideoDuration, percentWatched, resumeTime)

	// Only save if there was meaningful activity (at least 1% watched)
	if percentWatched < 1 {
		logger.Debugf("[DLNA Activity] Session too short, skipping save")
		return
	}

	// Skip DB operations if txnManager is nil (for testing)
	if t.txnManager == nil {
		logger.Debugf("[DLNA Activity] No transaction manager, skipping DB save")
		return
	}

	// Determine what needs to be saved
	shouldSaveResume := resumeTime > 0
	shouldAddView := !session.PlayCountAdded && percentWatched >= float64(t.getMinimumPlayPercent())

	// Nothing to save
	if !shouldSaveResume && !shouldAddView {
		return
	}

	// Save everything in a single transaction
	ctx := context.Background()
	if err := txn.WithTxn(ctx, t.txnManager, func(ctx context.Context) error {
		// Save resume time only. DLNA clients buffer aggressively and don't report
		// playback position, so we can't accurately track play duration - saving
		// guesses would corrupt analytics. Resume time is still useful as a
		// "continue watching" hint even if imprecise.
		if shouldSaveResume {
			if _, err := t.sceneWriter.SaveActivity(ctx, session.SceneID, &resumeTime, nil); err != nil {
				return fmt.Errorf("save resume time: %w", err)
			}
		}

		// Increment play count (also updates last_played_at via view date)
		if shouldAddView {
			if _, err := t.sceneWriter.AddViews(ctx, session.SceneID, []time.Time{time.Now()}); err != nil {
				return fmt.Errorf("add view: %w", err)
			}
			session.PlayCountAdded = true
			logger.Debugf("[DLNA Activity] Incremented play count for scene %d (%.1f%% watched)",
				session.SceneID, percentWatched)
		}

		return nil
	}); err != nil {
		logger.Warnf("[DLNA Activity] Failed to save activity for scene %d: %v", session.SceneID, err)
	}
}

// isEnabled returns true if activity tracking is enabled.
func (t *ActivityTracker) isEnabled() bool {
	if t.config == nil {
		return true // Default to enabled
	}
	return t.config.GetDLNAActivityTrackingEnabled()
}

// getMinimumPlayPercent returns the minimum play percentage for incrementing play count.
func (t *ActivityTracker) getMinimumPlayPercent() int {
	if t.config == nil {
		return 0 // Default: any play increments count (matches frontend default)
	}
	return t.config.GetMinimumPlayPercent()
}

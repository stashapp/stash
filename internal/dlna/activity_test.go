package dlna

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockSceneWriter is a mock implementation of SceneActivityWriter
type mockSceneWriter struct {
	mu                sync.Mutex
	saveActivityCalls []saveActivityCall
	addViewsCalls     []addViewsCall
}

type saveActivityCall struct {
	sceneID      int
	resumeTime   *float64
	playDuration *float64
}

type addViewsCall struct {
	sceneID int
	dates   []time.Time
}

func (m *mockSceneWriter) SaveActivity(_ context.Context, sceneID int, resumeTime *float64, playDuration *float64) (bool, error) {
	m.mu.Lock()
	m.saveActivityCalls = append(m.saveActivityCalls, saveActivityCall{
		sceneID:      sceneID,
		resumeTime:   resumeTime,
		playDuration: playDuration,
	})
	m.mu.Unlock()
	return true, nil
}

func (m *mockSceneWriter) AddViews(_ context.Context, sceneID int, dates []time.Time) ([]time.Time, error) {
	m.mu.Lock()
	m.addViewsCalls = append(m.addViewsCalls, addViewsCall{
		sceneID: sceneID,
		dates:   dates,
	})
	m.mu.Unlock()
	return dates, nil
}

// mockConfig is a mock implementation of ActivityConfig
type mockConfig struct {
	enabled        bool
	minPlayPercent int
}

func (c *mockConfig) GetDLNAActivityTrackingEnabled() bool {
	return c.enabled
}

func (c *mockConfig) GetMinimumPlayPercent() int {
	return c.minPlayPercent
}

func TestStreamSession_PercentWatched(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		startTime     time.Time
		lastActivity  time.Time
		videoDuration float64
		expected      float64
	}{
		{
			name:          "no video duration",
			startTime:     now.Add(-60 * time.Second),
			lastActivity:  now,
			videoDuration: 0,
			expected:      0,
		},
		{
			name:          "half watched",
			startTime:     now.Add(-60 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, watched for 1 minute = 50%
			expected:      50.0,
		},
		{
			name:          "fully watched",
			startTime:     now.Add(-120 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, watched for 2 minutes = 100%
			expected:      100.0,
		},
		{
			name:          "quarter watched",
			startTime:     now.Add(-30 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, watched for 30 seconds = 25%
			expected:      25.0,
		},
		{
			name:          "elapsed exceeds duration - capped at 100%",
			startTime:     now.Add(-180 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, but 3 minutes elapsed = capped at 100%
			expected:      100.0,
		},
		{
			name:          "no elapsed time",
			startTime:     now,
			lastActivity:  now,
			videoDuration: 120.0,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &streamSession{
				StartTime:     tt.startTime,
				LastActivity:  tt.lastActivity,
				VideoDuration: tt.videoDuration,
			}
			result := session.percentWatched()
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}

func TestStreamSession_EstimatedPlayDuration(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		startTime     time.Time
		lastActivity  time.Time
		videoDuration float64
		expected      float64
	}{
		{
			name:          "elapsed less than duration",
			startTime:     now.Add(-30 * time.Second),
			lastActivity:  now,
			videoDuration: 120,
			expected:      30.0,
		},
		{
			name:          "elapsed exceeds duration - capped",
			startTime:     now.Add(-180 * time.Second),
			lastActivity:  now,
			videoDuration: 120,
			expected:      120.0,
		},
		{
			name:          "no duration limit",
			startTime:     now.Add(-300 * time.Second),
			lastActivity:  now,
			videoDuration: 0,
			expected:      300.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &streamSession{
				StartTime:     tt.startTime,
				LastActivity:  tt.lastActivity,
				VideoDuration: tt.videoDuration,
			}
			result := session.estimatedPlayDuration()
			assert.InDelta(t, tt.expected, result, 1.0) // Allow 1 second tolerance
		})
	}
}

func TestStreamSession_EstimatedResumeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		startTime     time.Time
		lastActivity  time.Time
		videoDuration float64
		expected      float64
	}{
		{
			name:          "no elapsed time",
			startTime:     now,
			lastActivity:  now,
			videoDuration: 120.0,
			expected:      0,
		},
		{
			name:          "half way through",
			startTime:     now.Add(-60 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, watched for 1 minute = resume at 60s
			expected:      60.0,
		},
		{
			name:          "quarter way through",
			startTime:     now.Add(-30 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 2 minutes, watched for 30 seconds = resume at 30s
			expected:      30.0,
		},
		{
			name:          "98% complete - should reset to 0",
			startTime:     now.Add(-118 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 98.3% elapsed, should reset
			expected:      0,
		},
		{
			name:          "100% complete - should reset to 0",
			startTime:     now.Add(-120 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0,
			expected:      0,
		},
		{
			name:          "elapsed exceeds duration - capped and reset to 0",
			startTime:     now.Add(-180 * time.Second),
			lastActivity:  now,
			videoDuration: 120.0, // 150% elapsed, capped at 100%, reset to 0
			expected:      0,
		},
		{
			name:          "no video duration",
			startTime:     now.Add(-60 * time.Second),
			lastActivity:  now,
			videoDuration: 0,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &streamSession{
				StartTime:     tt.startTime,
				LastActivity:  tt.lastActivity,
				VideoDuration: tt.videoDuration,
			}
			result := session.estimatedResumeTime()
			assert.InDelta(t, tt.expected, result, 1.0) // Allow 1 second tolerance
		})
	}
}

func TestSessionKey(t *testing.T) {
	key := sessionKey("192.168.1.100", 42)
	assert.Equal(t, "192.168.1.100:42", key)
}

func TestActivityTracker_RecordRequest(t *testing.T) {
	config := &mockConfig{enabled: true, minPlayPercent: 50}

	// Create tracker without starting the goroutine (for unit testing)
	tracker := &ActivityTracker{
		txnManager:     nil, // Don't need DB for this test
		sceneWriter:    nil,
		config:         config,
		sessionTimeout: DefaultSessionTimeout,
		sessions:       make(map[string]*streamSession),
	}

	// Record first request - should create new session
	tracker.RecordRequest(42, "192.168.1.100", 120.0)

	tracker.mutex.RLock()
	session := tracker.sessions["192.168.1.100:42"]
	tracker.mutex.RUnlock()

	assert.NotNil(t, session)
	assert.Equal(t, 42, session.SceneID)
	assert.Equal(t, "192.168.1.100", session.ClientIP)
	assert.Equal(t, 120.0, session.VideoDuration)
	assert.False(t, session.StartTime.IsZero())
	assert.False(t, session.LastActivity.IsZero())

	// Record second request - should update LastActivity
	firstActivity := session.LastActivity
	time.Sleep(10 * time.Millisecond)
	tracker.RecordRequest(42, "192.168.1.100", 120.0)

	tracker.mutex.RLock()
	session = tracker.sessions["192.168.1.100:42"]
	tracker.mutex.RUnlock()

	assert.True(t, session.LastActivity.After(firstActivity))
}

func TestActivityTracker_DisabledTracking(t *testing.T) {
	config := &mockConfig{enabled: false, minPlayPercent: 50}

	// Create tracker without starting the goroutine (for unit testing)
	tracker := &ActivityTracker{
		txnManager:     nil,
		sceneWriter:    nil,
		config:         config,
		sessionTimeout: DefaultSessionTimeout,
		sessions:       make(map[string]*streamSession),
	}

	// Record request - should be ignored when tracking is disabled
	tracker.RecordRequest(42, "192.168.1.100", 120.0)

	tracker.mutex.RLock()
	sessionCount := len(tracker.sessions)
	tracker.mutex.RUnlock()

	assert.Equal(t, 0, sessionCount)
}

func TestActivityTracker_SessionExpiration(t *testing.T) {
	// For this test, we'll test the session expiration logic directly
	// without the full transaction manager integration

	sceneWriter := &mockSceneWriter{}
	config := &mockConfig{enabled: true, minPlayPercent: 10}

	// Create a tracker with nil txnManager - we'll test processCompletedSession separately
	// Here we just verify the session management logic
	tracker := &ActivityTracker{
		txnManager:     nil, // Skip DB calls for this test
		sceneWriter:    sceneWriter,
		config:         config,
		sessionTimeout: 100 * time.Millisecond,
		sessions:       make(map[string]*streamSession),
	}

	// Manually add a session
	// Use a short video duration (1 second) so the test can verify expiration quickly.
	now := time.Now()
	tracker.sessions["192.168.1.100:42"] = &streamSession{
		SceneID:       42,
		ClientIP:      "192.168.1.100",
		StartTime:     now.Add(-5 * time.Second),        // Started 5 seconds ago
		LastActivity:  now.Add(-200 * time.Millisecond), // Last activity 200ms ago (> 100ms timeout)
		VideoDuration: 1.0,                              // Short video so timeSinceStart > videoDuration
	}

	// Verify session exists
	assert.Len(t, tracker.sessions, 1)

	// Process expired sessions - this will try to save activity but txnManager is nil
	// so it will skip the DB calls but still remove the session
	tracker.processExpiredSessions()

	// Verify session was removed (even though DB calls were skipped)
	assert.Len(t, tracker.sessions, 0)
}

func TestActivityTracker_SessionExpiration_StoppedEarly(t *testing.T) {
	// Test that sessions expire when user stops watching early (before video ends)
	// This was a bug where sessions wouldn't expire until video duration passed

	config := &mockConfig{enabled: true, minPlayPercent: 10}
	tracker := &ActivityTracker{
		txnManager:     nil,
		sceneWriter:    nil,
		config:         config,
		sessionTimeout: 100 * time.Millisecond,
		sessions:       make(map[string]*streamSession),
	}

	// User started watching a 30-minute video but stopped after 5 seconds
	now := time.Now()
	tracker.sessions["192.168.1.100:42"] = &streamSession{
		SceneID:       42,
		ClientIP:      "192.168.1.100",
		StartTime:     now.Add(-5 * time.Second),        // Started 5 seconds ago
		LastActivity:  now.Add(-200 * time.Millisecond), // Last activity 200ms ago (> 100ms timeout)
		VideoDuration: 1800.0,                           // 30 minute video - much longer than elapsed time
	}

	assert.Len(t, tracker.sessions, 1)

	// Session should expire because timeSinceActivity > timeout
	// Even though the video is 30 minutes and only 5 seconds have passed
	tracker.processExpiredSessions()

	// Verify session was expired
	assert.Len(t, tracker.sessions, 0, "Session should expire when user stops early, not wait for video duration")
}

func TestActivityTracker_MinimumPlayPercentThreshold(t *testing.T) {
	// Test the threshold logic without full transaction integration
	config := &mockConfig{enabled: true, minPlayPercent: 75} // High threshold

	tracker := &ActivityTracker{
		txnManager:     nil,
		sceneWriter:    nil,
		config:         config,
		sessionTimeout: 50 * time.Millisecond,
		sessions:       make(map[string]*streamSession),
	}

	// Test that getMinimumPlayPercent returns the configured value
	assert.Equal(t, 75, tracker.getMinimumPlayPercent())

	// Create a session with 30% watched (36 seconds of a 120 second video)
	now := time.Now()
	session := &streamSession{
		SceneID:       42,
		StartTime:     now.Add(-36 * time.Second),
		LastActivity:  now,
		VideoDuration: 120.0,
	}

	// 30% is below 75% threshold
	percentWatched := session.percentWatched()
	assert.InDelta(t, 30.0, percentWatched, 0.1)
	assert.False(t, percentWatched >= float64(tracker.getMinimumPlayPercent()))
}

func TestActivityTracker_MultipleSessions(t *testing.T) {
	config := &mockConfig{enabled: true, minPlayPercent: 50}

	// Create tracker without starting the goroutine (for unit testing)
	tracker := &ActivityTracker{
		txnManager:     nil,
		sceneWriter:    nil,
		config:         config,
		sessionTimeout: DefaultSessionTimeout,
		sessions:       make(map[string]*streamSession),
	}

	// Different clients watching same scene
	tracker.RecordRequest(42, "192.168.1.100", 120.0)
	tracker.RecordRequest(42, "192.168.1.101", 120.0)

	// Same client watching different scenes
	tracker.RecordRequest(43, "192.168.1.100", 180.0)

	tracker.mutex.RLock()
	assert.Len(t, tracker.sessions, 3)
	tracker.mutex.RUnlock()
}

func TestActivityTracker_ShortSessionIgnored(t *testing.T) {
	// Test that short sessions are ignored
	// Create a session with only ~0.8% watched (1 second of a 120 second video)
	now := time.Now()
	session := &streamSession{
		SceneID:       42,
		ClientIP:      "192.168.1.100",
		StartTime:     now.Add(-1 * time.Second), // Only 1 second
		LastActivity:  now,
		VideoDuration: 120.0, // 2 minutes
	}

	// Verify percent watched is below threshold (1s / 120s = 0.83%)
	assert.InDelta(t, 0.83, session.percentWatched(), 0.1)

	// Verify play duration is short
	assert.InDelta(t, 1.0, session.estimatedPlayDuration(), 0.5)

	// Both are below the minimum thresholds (1% and 5 seconds)
	percentWatched := session.percentWatched()
	playDuration := session.estimatedPlayDuration()
	shouldSkip := percentWatched < 1 && playDuration < 5
	assert.True(t, shouldSkip, "Short session should be skipped")
}

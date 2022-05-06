package astikit

import (
	"context"
	"sync"
	"time"
)

// Limiter represents a limiter
type Limiter struct {
	buckets map[string]*LimiterBucket
	m       *sync.Mutex // Locks buckets
}

// NewLimiter creates a new limiter
func NewLimiter() *Limiter {
	return &Limiter{
		buckets: make(map[string]*LimiterBucket),
		m:       &sync.Mutex{},
	}
}

// Add adds a new bucket
func (l *Limiter) Add(name string, cap int, period time.Duration) *LimiterBucket {
	l.m.Lock()
	defer l.m.Unlock()
	if _, ok := l.buckets[name]; !ok {
		l.buckets[name] = newLimiterBucket(cap, period)
	}
	return l.buckets[name]
}

// Bucket retrieves a bucket from the limiter
func (l *Limiter) Bucket(name string) (b *LimiterBucket, ok bool) {
	l.m.Lock()
	defer l.m.Unlock()
	b, ok = l.buckets[name]
	return
}

// Close closes the limiter properly
func (l *Limiter) Close() {
	l.m.Lock()
	defer l.m.Unlock()
	for _, b := range l.buckets {
		b.Close()
	}
}

// LimiterBucket represents a limiter bucket
type LimiterBucket struct {
	cancel context.CancelFunc
	cap    int
	ctx    context.Context
	count  int
	period time.Duration
	o      *sync.Once
}

// newLimiterBucket creates a new bucket
func newLimiterBucket(cap int, period time.Duration) (b *LimiterBucket) {
	b = &LimiterBucket{
		cap:    cap,
		count:  0,
		period: period,
		o:      &sync.Once{},
	}
	b.ctx, b.cancel = context.WithCancel(context.Background())
	go b.tick()
	return
}

// Inc increments the bucket count
func (b *LimiterBucket) Inc() bool {
	if b.count >= b.cap {
		return false
	}
	b.count++
	return true
}

// tick runs a ticker to purge the bucket
func (b *LimiterBucket) tick() {
	var t = time.NewTicker(b.period)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			b.count = 0
		case <-b.ctx.Done():
			return
		}
	}
}

// close closes the bucket properly
func (b *LimiterBucket) Close() {
	b.o.Do(func() {
		b.cancel()
	})
}

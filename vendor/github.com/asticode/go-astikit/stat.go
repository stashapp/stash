package astikit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Stater is an object that can compute and handle stats
type Stater struct {
	cancel  context.CancelFunc
	ctx     context.Context
	h       StatsHandleFunc
	m       *sync.Mutex // Locks ss
	period  time.Duration
	running uint32
	ss      map[*StatMetadata]StatOptions
}

// StatOptions represents stat options
type StatOptions struct {
	Metadata *StatMetadata
	// Either a StatValuer or StatValuerOverTime
	Valuer interface{}
}

// StatsHandleFunc is a method that can handle stat values
type StatsHandleFunc func(stats []StatValue)

// StatMetadata represents a stat metadata
type StatMetadata struct {
	Description string
	Label       string
	Name        string
	Unit        string
}

// StatValuer represents a stat valuer
type StatValuer interface {
	Value() interface{}
}

// StatValuerOverTime represents a stat valuer over time
type StatValuerOverTime interface {
	Value(delta time.Duration) interface{}
}

// StatValue represents a stat value
type StatValue struct {
	*StatMetadata
	Value interface{}
}

// StaterOptions represents stater options
type StaterOptions struct {
	HandleFunc StatsHandleFunc
	Period     time.Duration
}

// NewStater creates a new stater
func NewStater(o StaterOptions) *Stater {
	return &Stater{
		h:      o.HandleFunc,
		m:      &sync.Mutex{},
		period: o.Period,
		ss:     make(map[*StatMetadata]StatOptions),
	}
}

// Start starts the stater
func (s *Stater) Start(ctx context.Context) {
	// Check context
	if ctx.Err() != nil {
		return
	}

	// Make sure to start only once
	if atomic.CompareAndSwapUint32(&s.running, 0, 1) {
		// Update status
		defer atomic.StoreUint32(&s.running, 0)

		// Reset context
		s.ctx, s.cancel = context.WithCancel(ctx)

		// Create ticker
		t := time.NewTicker(s.period)
		defer t.Stop()

		// Loop
		lastStatAt := now()
		for {
			select {
			case <-t.C:
				// Get delta
				n := now()
				delta := n.Sub(lastStatAt)
				lastStatAt = n

				// Loop through stats
				var stats []StatValue
				s.m.Lock()
				for _, o := range s.ss {
					// Get value
					var v interface{}
					if h, ok := o.Valuer.(StatValuer); ok {
						v = h.Value()
					} else if h, ok := o.Valuer.(StatValuerOverTime); ok {
						v = h.Value(delta)
					} else {
						continue
					}

					// Append
					stats = append(stats, StatValue{
						StatMetadata: o.Metadata,
						Value:        v,
					})
				}
				s.m.Unlock()

				// Handle stats
				go s.h(stats)
			case <-s.ctx.Done():
				return
			}
		}
	}
}

// Stop stops the stater
func (s *Stater) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// AddStats adds stats
func (s *Stater) AddStats(os ...StatOptions) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, o := range os {
		s.ss[o.Metadata] = o
	}
}

// DelStats deletes stats
func (s *Stater) DelStats(os ...StatOptions) {
	s.m.Lock()
	defer s.m.Unlock()
	for _, o := range os {
		delete(s.ss, o.Metadata)
	}
}

type durationStatOverTime struct {
	d           time.Duration
	fn          func(d, delta time.Duration) interface{}
	m           *sync.Mutex // Locks isStarted
	lastBeginAt time.Time
}

func newDurationStatOverTime(fn func(d, delta time.Duration) interface{}) *durationStatOverTime {
	return &durationStatOverTime{
		fn: fn,
		m:  &sync.Mutex{},
	}
}

func (s *durationStatOverTime) Begin() {
	s.m.Lock()
	defer s.m.Unlock()
	s.lastBeginAt = now()
}

func (s *durationStatOverTime) End() {
	s.m.Lock()
	defer s.m.Unlock()
	s.d += now().Sub(s.lastBeginAt)
	s.lastBeginAt = time.Time{}
}

func (s *durationStatOverTime) Value(delta time.Duration) (o interface{}) {
	// Lock
	s.m.Lock()
	defer s.m.Unlock()

	// Get current values
	n := now()
	d := s.d

	// Recording is still in process
	if !s.lastBeginAt.IsZero() {
		d += n.Sub(s.lastBeginAt)
		s.lastBeginAt = n
	}

	// Compute stat
	o = s.fn(d, delta)
	s.d = 0
	return
}

// DurationPercentageStat is an object capable of computing the percentage of time some work is taking per second
type DurationPercentageStat struct {
	*durationStatOverTime
}

// NewDurationPercentageStat creates a new duration percentage stat
func NewDurationPercentageStat() *DurationPercentageStat {
	return &DurationPercentageStat{durationStatOverTime: newDurationStatOverTime(func(d, delta time.Duration) interface{} {
		if delta == 0 {
			return 0
		}
		return float64(d) / float64(delta) * 100
	})}
}

type counterStatOverTime struct {
	c  float64
	fn func(c, t float64, delta time.Duration) interface{}
	m  *sync.Mutex // Locks isStarted
	t  float64
}

func newCounterStatOverTime(fn func(c, t float64, delta time.Duration) interface{}) *counterStatOverTime {
	return &counterStatOverTime{
		fn: fn,
		m:  &sync.Mutex{},
	}
}

func (s *counterStatOverTime) Add(delta float64) {
	s.m.Lock()
	defer s.m.Unlock()
	s.c += delta
	s.t++
}

func (s *counterStatOverTime) Value(delta time.Duration) interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	c := s.c
	t := s.t
	s.c = 0
	s.t = 0
	return s.fn(c, t, delta)
}

// CounterAvgStat is an object capable of computing the average value of a counter
type CounterAvgStat struct {
	*counterStatOverTime
}

// NewCounterAvgStat creates a new counter avg stat
func NewCounterAvgStat() *CounterAvgStat {
	return &CounterAvgStat{counterStatOverTime: newCounterStatOverTime(func(c, t float64, delta time.Duration) interface{} {
		if t == 0 {
			return 0
		}
		return c / t
	})}
}

// CounterRateStat is an object capable of computing the average value of a counter per second
type CounterRateStat struct {
	*counterStatOverTime
}

// NewCounterRateStat creates a new counter rate stat
func NewCounterRateStat() *CounterRateStat {
	return &CounterRateStat{counterStatOverTime: newCounterStatOverTime(func(c, t float64, delta time.Duration) interface{} {
		if delta.Seconds() == 0 {
			return 0
		}
		return c / delta.Seconds()
	})}
}

// CounterStat is an object capable of computing a counter that never gets reset
type CounterStat struct {
	c float64
	m *sync.Mutex
}

// NewCounterStat creates a new counter stat
func NewCounterStat() *CounterStat {
	return &CounterStat{m: &sync.Mutex{}}
}

func (s *CounterStat) Add(delta float64) {
	s.m.Lock()
	defer s.m.Unlock()
	s.c += delta
}

func (s *CounterStat) Value() interface{} {
	s.m.Lock()
	defer s.m.Unlock()
	return s.c
}

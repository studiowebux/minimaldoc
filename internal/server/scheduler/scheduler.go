// Package scheduler provides background job scheduling for minimaldoc-server.
package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/studiowebux/minimaldoc/internal/server/store"
)

// Scheduler runs background jobs on a ticker.
type Scheduler struct {
	db       store.Store
	interval time.Duration
	ticker   *time.Ticker
	done     chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// New creates a new scheduler with the given interval.
func New(db store.Store, interval time.Duration) *Scheduler {
	if interval < time.Second {
		interval = time.Minute // Default to 1 minute
	}
	return &Scheduler{
		db:       db,
		interval: interval,
		done:     make(chan struct{}),
	}
}

// Start begins the scheduler loop.
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.ticker = time.NewTicker(s.interval)
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run()

	slog.Info("scheduler started", "interval", s.interval)
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.done)
	s.ticker.Stop()
	s.wg.Wait()

	slog.Info("scheduler stopped")
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	// Run immediately on start
	s.runJobs()

	for {
		select {
		case <-s.done:
			return
		case <-s.ticker.C:
			s.runJobs()
		}
	}
}

func (s *Scheduler) runJobs() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Publish scheduled posts
	s.publishScheduledPosts(ctx)
}

func (s *Scheduler) publishScheduledPosts(ctx context.Context) {
	count, err := s.db.PublishScheduledPosts(ctx)
	if err != nil {
		slog.Error("failed to publish scheduled posts", "error", err)
		return
	}
	if count > 0 {
		slog.Info("published scheduled posts", "count", count)
	}
}

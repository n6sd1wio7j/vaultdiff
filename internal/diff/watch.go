package diff

import (
	"context"
	"time"
)

// WatchOptions configures the polling behavior.
type WatchOptions struct {
	Interval  time.Duration
	MaxChecks int // 0 means unlimited
}

// WatchResult holds a single poll result.
type WatchResult struct {
	CheckedAt time.Time
	Entries   []Entry
	HasDrift  bool
}

// Watcher polls two secret versions and emits results when drift is detected.
type Watcher struct {
	opts WatchOptions
	fetch func() ([]Entry, error)
}

// NewWatcher creates a Watcher with the given fetch function and options.
func NewWatcher(fetch func() ([]Entry, error), opts WatchOptions) *Watcher {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}
	return &Watcher{opts: opts, fetch: fetch}
}

// Run starts polling and sends results to the returned channel.
// The channel is closed when ctx is cancelled or MaxChecks is reached.
func (w *Watcher) Run(ctx context.Context) <-chan WatchResult {
	ch := make(chan WatchResult)
	go func() {
		defer close(ch)
		checks := 0
		for {
			entries, err := w.fetch()
			if err == nil {
				result := WatchResult{
					CheckedAt: time.Now().UTC(),
					Entries:   entries,
					HasDrift:  HasChanges(entries),
				}
				select {
				case ch <- result:
				case <-ctx.Done():
					return
				}
			}
			checks++
			if w.opts.MaxChecks > 0 && checks >= w.opts.MaxChecks {
				return
			}
			select {
			case <-time.After(w.opts.Interval):
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

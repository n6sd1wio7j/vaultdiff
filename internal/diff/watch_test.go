package diff

import (
	"context"
	"testing"
	"time"
)

func fixedEntries(entries []Entry) func() ([]Entry, error) {
	return func() ([]Entry, error) { return entries, nil }
}

func TestWatcher_EmitsResults(t *testing.T) {
	entries := []Entry{
		{Key: "foo", OldValue: "a", NewValue: "b", Status: StatusModified},
	}
	w := NewWatcher(fixedEntries(entries), WatchOptions{Interval: 1 * time.Millisecond, MaxChecks: 2})
	ctx := context.Background()
	ch := w.Run(ctx)
	count := 0
	for r := range ch {
		count++
		if !r.HasDrift {
			t.Error("expected HasDrift true")
		}
	}
	if count != 2 {
		t.Errorf("expected 2 results, got %d", count)
	}
}

func TestWatcher_NoDrift(t *testing.T) {
	entries := []Entry{
		{Key: "foo", OldValue: "a", NewValue: "a", Status: StatusUnchanged},
	}
	w := NewWatcher(fixedEntries(entries), WatchOptions{Interval: 1 * time.Millisecond, MaxChecks: 1})
	for r := range w.Run(context.Background()) {
		if r.HasDrift {
			t.Error("expected no drift")
		}
	}
}

func TestWatcher_CancelStops(t *testing.T) {
	w := NewWatcher(fixedEntries(nil), WatchOptions{Interval: 10 * time.Second})
	ctx, cancel := context.WithCancel(context.Background())
	ch := w.Run(ctx)
	<-ch // first result
	cancel()
	_, open := <-ch
	if open {
		t.Error("expected channel closed after cancel")
	}
}

func TestWatcher_DefaultInterval(t *testing.T) {
	w := NewWatcher(fixedEntries(nil), WatchOptions{})
	if w.opts.Interval != 30*time.Second {
		t.Errorf("expected 30s default, got %v", w.opts.Interval)
	}
}

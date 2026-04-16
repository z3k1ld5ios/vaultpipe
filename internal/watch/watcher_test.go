package watch_test

import (
	"context"
	"testing"
	"time"

	"github.com/your-org/vaultpipe/internal/watch"
)

func staticFetcher(secrets map[string]string) watch.SecretFetcher {
	return func(_ context.Context, _ string) (map[string]string, error) {
		return secrets, nil
	}
}

func TestWatch_EmitsOnChange(t *testing.T) {
	calls := 0
	secretSets := []map[string]string{
		{"KEY": "v1"},
		{"KEY": "v2"},
	}
	fetcher := func(_ context.Context, _ string) (map[string]string, error) {
		if calls >= len(secretSets) {
			return secretSets[len(secretSets)-1], nil
		}
		s := secretSets[calls]
		calls++
		return s, nil
	}
	w := watch.NewWatcher(fetcher, "secret/app", 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	ch := w.Watch(ctx)
	var events []watch.ChangeEvent
	for e := range ch {
		events = append(events, e)
	}
	if len(events) < 2 {
		t.Fatalf("expected at least 2 change events, got %d", len(events))
	}
	if events[0].Secrets["KEY"] != "v1" {
		t.Errorf("expected v1, got %s", events[0].Secrets["KEY"])
	}
	if events[1].Secrets["KEY"] != "v2" {
		t.Errorf("expected v2, got %s", events[1].Secrets["KEY"])
	}
}

func TestWatch_NoEmitWhenUnchanged(t *testing.T) {
	w := watch.NewWatcher(staticFetcher(map[string]string{"K": "same"}), "secret/app", 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	ch := w.Watch(ctx)
	count := 0
	for range ch {
		count++
	}
	if count != 1 {
		t.Errorf("expected exactly 1 event (initial change), got %d", count)
	}
}

func TestWatch_ContextCancel_ClosesChannel(t *testing.T) {
	w := watch.NewWatcher(staticFetcher(map[string]string{}), "secret/app", 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	ch := w.Watch(ctx)
	cancel()
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("channel was not closed after context cancel")
	}
}

package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
)

func sampleSecrets() map[string]string {
	return map[string]string{"API_KEY": "abc123", "DB_PASS": "secret"}
}

func TestGet_MissOnEmpty(t *testing.T) {
	c := cache.New(5 * time.Minute)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss on empty cache")
	}
}

func TestSetAndGet_Hit(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", sampleSecrets())
	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("unexpected value: %s", got["API_KEY"])
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("secret/app", sampleSecrets())
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestSet_ZeroTTL_NeverStores(t *testing.T) {
	c := cache.New(0)
	c.Set("secret/app", sampleSecrets())
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss when TTL is zero")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", sampleSecrets())
	c.Invalidate("secret/app")
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", sampleSecrets())
	c.Set("secret/db", map[string]string{"PASS": "x"})
	c.Flush()
	if _, ok := c.Get("secret/app"); ok {
		t.Error("expected miss for secret/app after flush")
	}
	if _, ok := c.Get("secret/db"); ok {
		t.Error("expected miss for secret/db after flush")
	}
}

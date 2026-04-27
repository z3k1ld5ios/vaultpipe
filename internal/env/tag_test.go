package env

import (
	"testing"
)

func TestSet_AddsTag(t *testing.T) {
	tr := NewTagger()
	tr.Set("DB_HOST", "source", "vault")
	tags := tr.Get("DB_HOST")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Key != "source" || tags[0].Value != "vault" {
		t.Errorf("unexpected tag: %+v", tags[0])
	}
}

func TestSet_OverwritesDuplicateTagKey(t *testing.T) {
	tr := NewTagger()
	tr.Set("API_KEY", "source", "vault")
	tr.Set("API_KEY", "source", "env")
	tags := tr.Get("API_KEY")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag after overwrite, got %d", len(tags))
	}
	if tags[0].Value != "env" {
		t.Errorf("expected value 'env', got %q", tags[0].Value)
	}
}

func TestGet_CaseInsensitiveKey(t *testing.T) {
	tr := NewTagger()
	tr.Set("db_host", "tier", "prod")
	tags := tr.Get("DB_HOST")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
}

func TestGet_MissingKey_ReturnsNil(t *testing.T) {
	tr := NewTagger()
	if tags := tr.Get("NONEXISTENT"); tags != nil {
		t.Errorf("expected nil, got %v", tags)
	}
}

func TestHasTag_Match(t *testing.T) {
	tr := NewTagger()
	tr.Set("SECRET_KEY", "sensitive", "true")
	if !tr.HasTag("SECRET_KEY", "sensitive", "true") {
		t.Error("expected HasTag to return true")
	}
}

func TestHasTag_NoMatch(t *testing.T) {
	tr := NewTagger()
	tr.Set("SECRET_KEY", "sensitive", "true")
	if tr.HasTag("SECRET_KEY", "sensitive", "false") {
		t.Error("expected HasTag to return false")
	}
}

func TestFilter_RetainsMatchingKeys(t *testing.T) {
	tr := NewTagger()
	tr.Set("TOKEN", "sensitive", "true")
	tr.Set("HOST", "sensitive", "false")
	tr.Set("PORT", "tier", "prod")

	env := map[string]string{"TOKEN": "abc", "HOST": "localhost", "PORT": "8080"}
	result := tr.Filter(env, "sensitive", "true")

	if len(result) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result))
	}
	if result["TOKEN"] != "abc" {
		t.Errorf("expected TOKEN=abc, got %q", result["TOKEN"])
	}
}

func TestFilter_EmptyEnv_ReturnsEmpty(t *testing.T) {
	tr := NewTagger()
	result := tr.Filter(map[string]string{}, "source", "vault")
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestSummary_WithTags(t *testing.T) {
	tr := NewTagger()
	tr.Set("DB_PASS", "source", "vault")
	tr.Set("DB_PASS", "sensitive", "true")
	s := tr.Summary("DB_PASS")
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_NoTags(t *testing.T) {
	tr := NewTagger()
	s := tr.Summary("UNKNOWN")
	expected := "UNKNOWN: (no tags)"
	if s != expected {
		t.Errorf("expected %q, got %q", expected, s)
	}
}

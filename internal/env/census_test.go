package env

import (
	"testing"
)

func TestAnalyze_BasicCounts(t *testing.T) {
	m := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"PASS": "",
	}

	c := Analyze(m)

	if c.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", c.TotalKeys)
	}
	if c.EmptyValues != 1 {
		t.Errorf("expected EmptyValues=1, got %d", c.EmptyValues)
	}
	if len(c.UniqueKeys) != 3 {
		t.Errorf("expected 3 unique keys, got %d", len(c.UniqueKeys))
	}
}

func TestAnalyze_EmptyMap(t *testing.T) {
	c := Analyze(map[string]string{})

	if c.TotalKeys != 0 {
		t.Errorf("expected TotalKeys=0, got %d", c.TotalKeys)
	}
	if c.EmptyValues != 0 {
		t.Errorf("expected EmptyValues=0, got %d", c.EmptyValues)
	}
}

func TestAnalyzeMultiple_DetectsDuplicates(t *testing.T) {
	src1 := map[string]string{"HOST": "a", "PORT": "1"}
	src2 := map[string]string{"HOST": "b", "TOKEN": "xyz"}

	c := AnalyzeMultiple(src1, src2)

	if c.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", c.TotalKeys)
	}
	if len(c.DuplicateKeys) != 1 || c.DuplicateKeys[0] != "HOST" {
		t.Errorf("expected DuplicateKeys=[HOST], got %v", c.DuplicateKeys)
	}
}

func TestAnalyzeMultiple_NoDuplicates(t *testing.T) {
	src1 := map[string]string{"A": "1"}
	src2 := map[string]string{"B": "2"}

	c := AnalyzeMultiple(src1, src2)

	if len(c.DuplicateKeys) != 0 {
		t.Errorf("expected no duplicates, got %v", c.DuplicateKeys)
	}
	if len(c.UniqueKeys) != 2 {
		t.Errorf("expected 2 unique keys, got %d", len(c.UniqueKeys))
	}
}

func TestAnalyzeMultiple_EmptySources(t *testing.T) {
	c := AnalyzeMultiple()

	if c.TotalKeys != 0 {
		t.Errorf("expected TotalKeys=0, got %d", c.TotalKeys)
	}
}

func TestCensus_Summary(t *testing.T) {
	c := Census{
		TotalKeys:     5,
		EmptyValues:   1,
		DuplicateKeys: []string{"HOST", "PORT"},
	}

	got := c.Summary()
	want := "total=5 empty=1 duplicates=2"

	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

package diff

import (
	"testing"
)

func makeMap(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCompare_AllMatch(t *testing.T) {
	base := makeMap("KEY", "value", "PORT", "8080")
	target := makeMap("KEY", "value", "PORT", "8080")
	results := Compare(base, target)
	for _, r := range results {
		if r.Status != StatusMatch {
			t.Errorf("expected match for key %q, got %s", r.Key, r.Status)
		}
	}
}

func TestCompare_MissingKey(t *testing.T) {
	base := makeMap("KEY", "value", "SECRET", "abc")
	target := makeMap("KEY", "value")
	results := Compare(base, target)
	missing := FilterByStatus(results, StatusMissing)
	if len(missing) != 1 || missing[0].Key != "SECRET" {
		t.Errorf("expected SECRET to be missing, got %+v", missing)
	}
}

func TestCompare_ExtraKey(t *testing.T) {
	base := makeMap("KEY", "value")
	target := makeMap("KEY", "value", "EXTRA", "bonus")
	results := Compare(base, target)
	extras := FilterByStatus(results, StatusExtra)
	if len(extras) != 1 || extras[0].Key != "EXTRA" {
		t.Errorf("expected EXTRA to be extra, got %+v", extras)
	}
}

func TestCompare_MismatchValue(t *testing.T) {
	base := makeMap("DB_HOST", "localhost")
	target := makeMap("DB_HOST", "prod.db.example.com")
	results := Compare(base, target)
	mismatches := FilterByStatus(results, StatusMismatch)
	if len(mismatches) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(mismatches))
	}
	if mismatches[0].BaseValue != "localhost" || mismatches[0].TargetValue != "prod.db.example.com" {
		t.Errorf("unexpected mismatch values: %+v", mismatches[0])
	}
}

func TestFilterByStatus(t *testing.T) {
	results := []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMissing},
		{Key: "C", Status: StatusExtra},
		{Key: "D", Status: StatusMismatch},
	}
	got := FilterByStatus(results, StatusMissing, StatusExtra)
	if len(got) != 2 {
		t.Errorf("expected 2 results, got %d", len(got))
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	results := Compare(map[string]string{}, map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected no results for empty maps, got %d", len(results))
	}
}

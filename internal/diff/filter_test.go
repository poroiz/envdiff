package diff

import (
	"testing"
)

func makeFilterEntries() []Entry {
	return []Entry{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMissing},
		{Key: "C", Status: StatusExtra},
		{Key: "D", Status: StatusMismatch},
	}
}

func TestDefaultFilter_ExcludesMatch(t *testing.T) {
	f := DefaultFilter()
	result := f.Apply(makeFilterEntries())
	for _, e := range result {
		if e.Status == StatusMatch {
			t.Errorf("DefaultFilter should exclude match entries, got key %q", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestAllFilter_IncludesAll(t *testing.T) {
	f := AllFilter()
	result := f.Apply(makeFilterEntries())
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	f := StatusFilter{IncludeMissing: true}
	result := f.Apply(makeFilterEntries())
	if len(result) != 1 || result[0].Key != "B" {
		t.Errorf("expected only missing entry B, got %+v", result)
	}
}

func TestParseStatusFilter_Empty(t *testing.T) {
	f, err := ParseStatusFilter([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IncludeMissing || !f.IncludeExtra || !f.IncludeMismatch {
		t.Error("empty input should return DefaultFilter with missing/extra/mismatch enabled")
	}
	if f.IncludeMatch {
		t.Error("DefaultFilter should not include match")
	}
}

func TestParseStatusFilter_ValidValues(t *testing.T) {
	f, err := ParseStatusFilter([]string{"missing", "match"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IncludeMissing {
		t.Error("expected IncludeMissing to be true")
	}
	if !f.IncludeMatch {
		t.Error("expected IncludeMatch to be true")
	}
	if f.IncludeExtra || f.IncludeMismatch {
		t.Error("expected IncludeExtra and IncludeMismatch to be false")
	}
}

func TestParseStatusFilter_InvalidValue(t *testing.T) {
	_, err := ParseStatusFilter([]string{"missing", "unknown"})
	if err == nil {
		t.Error("expected error for unknown status filter value")
	}
}

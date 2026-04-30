package diff

import (
	"testing"
)

func makeSortEntries() []Entry {
	return []Entry{
		{Key: "ZEBRA", Status: StatusMatch},
		{Key: "ALPHA", Status: StatusMissing},
		{Key: "MONGO", Status: StatusMismatch},
		{Key: "BETA", Status: StatusExtra},
		{Key: "GAMMA", Status: StatusMissing},
	}
}

func TestSortEntries_ByKey(t *testing.T) {
	entries := makeSortEntries()
	sorted := SortEntries(entries, SortByKey)

	expectedKeys := []string{"ALPHA", "BETA", "GAMMA", "MONGO", "ZEBRA"}
	for i, key := range expectedKeys {
		if sorted[i].Key != key {
			t.Errorf("index %d: expected key %q, got %q", i, key, sorted[i].Key)
		}
	}
}

func TestSortEntries_ByStatus(t *testing.T) {
	entries := makeSortEntries()
	sorted := SortEntries(entries, SortByStatus)

	// Missing entries should come first, then extra, mismatch, match.
	if sorted[0].Status != StatusMissing || sorted[1].Status != StatusMissing {
		t.Errorf("expected first two entries to be Missing, got %v and %v", sorted[0].Status, sorted[1].Status)
	}
	if sorted[2].Status != StatusExtra {
		t.Errorf("expected index 2 to be Extra, got %v", sorted[2].Status)
	}
	if sorted[3].Status != StatusMismatch {
		t.Errorf("expected index 3 to be Mismatch, got %v", sorted[3].Status)
	}
	if sorted[4].Status != StatusMatch {
		t.Errorf("expected index 4 to be Match, got %v", sorted[4].Status)
	}
}

func TestSortEntries_DoesNotMutateOriginal(t *testing.T) {
	entries := makeSortEntries()
	originalFirst := entries[0].Key
	SortEntries(entries, SortByKey)
	if entries[0].Key != originalFirst {
		t.Errorf("original slice was mutated: expected %q at index 0, got %q", originalFirst, entries[0].Key)
	}
}

func TestParseSortOrder_Valid(t *testing.T) {
	if order, ok := ParseSortOrder("key"); !ok || order != SortByKey {
		t.Errorf("expected SortByKey for 'key'")
	}
	if order, ok := ParseSortOrder("status"); !ok || order != SortByStatus {
		t.Errorf("expected SortByStatus for 'status'")
	}
}

func TestParseSortOrder_Invalid(t *testing.T) {
	order, ok := ParseSortOrder("unknown")
	if ok {
		t.Errorf("expected ok=false for unknown sort order")
	}
	if order != SortByKey {
		t.Errorf("expected default SortByKey for unknown sort order")
	}
}

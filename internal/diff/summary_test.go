package diff

import (
	"strings"
	"testing"
)

func makeSummaryEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMatch},
		{Key: "C", Status: StatusMissing},
		{Key: "D", Status: StatusExtra},
		{Key: "E", Status: StatusMismatch},
	}
}

func TestSummarize_Counts(t *testing.T) {
	entries := makeSummaryEntries()
	s := Summarize(entries)

	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Counts.Match != 2 {
		t.Errorf("expected Match=2, got %d", s.Counts.Match)
	}
	if s.Counts.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", s.Counts.Missing)
	}
	if s.Counts.Extra != 1 {
		t.Errorf("expected Extra=1, got %d", s.Counts.Extra)
	}
	if s.Counts.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", s.Counts.Mismatch)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]DiffEntry{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
	if s.HasIssues() {
		t.Error("expected HasIssues=false for empty entries")
	}
}

func TestSummary_HasIssues_True(t *testing.T) {
	s := Summarize(makeSummaryEntries())
	if !s.HasIssues() {
		t.Error("expected HasIssues=true")
	}
}

func TestSummary_HasIssues_False(t *testing.T) {
	entries := []DiffEntry{
		{Key: "X", Status: StatusMatch},
		{Key: "Y", Status: StatusMatch},
	}
	s := Summarize(entries)
	if s.HasIssues() {
		t.Error("expected HasIssues=false when all keys match")
	}
}

func TestSummary_String(t *testing.T) {
	s := Summarize(makeSummaryEntries())
	out := s.String()
	for _, substr := range []string{"Total: 5", "Match: 2", "Missing: 1", "Extra: 1", "Mismatch: 1"} {
		if !strings.Contains(out, substr) {
			t.Errorf("expected summary string to contain %q, got: %s", substr, out)
		}
	}
}

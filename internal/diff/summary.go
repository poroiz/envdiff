package diff

import "fmt"

// StatusCount holds the count of each diff status.
type StatusCount struct {
	Match    int
	Missing  int
	Extra    int
	Mismatch int
}

// Summary aggregates diff results into a StatusCount and provides
// a human-readable overview string.
type Summary struct {
	Counts StatusCount
	Total  int
}

// Summarize computes a Summary from a slice of DiffEntry.
func Summarize(entries []DiffEntry) Summary {
	var counts StatusCount
	for _, e := range entries {
		switch e.Status {
		case StatusMatch:
			counts.Match++
		case StatusMissing:
			counts.Missing++
		case StatusExtra:
			counts.Extra++
		case StatusMismatch:
			counts.Mismatch++
		}
	}
	return Summary{
		Counts: counts,
		Total:  len(entries),
	}
}

// String returns a one-line summary suitable for CLI output.
func (s Summary) String() string {
	return fmt.Sprintf(
		"Total: %d | Match: %d | Missing: %d | Extra: %d | Mismatch: %d",
		s.Total, s.Counts.Match, s.Counts.Missing, s.Counts.Extra, s.Counts.Mismatch,
	)
}

// HasIssues returns true when there are any non-matching entries.
func (s Summary) HasIssues() bool {
	return s.Counts.Missing > 0 || s.Counts.Extra > 0 || s.Counts.Mismatch > 0
}

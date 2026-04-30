package diff

import "sort"

// SortOrder defines the ordering for diff entries.
type SortOrder int

const (
	// SortByKey sorts entries alphabetically by key name.
	SortByKey SortOrder = iota
	// SortByStatus sorts entries by status (missing, extra, mismatch, match).
	SortByStatus
)

// statusRank maps a Status to a numeric rank for ordering.
var statusRank = map[Status]int{
	StatusMissing:  0,
	StatusExtra:    1,
	StatusMismatch: 2,
	StatusMatch:    3,
}

// SortEntries returns a new sorted slice of Entry values based on the given SortOrder.
// The original slice is not modified.
func SortEntries(entries []Entry, order SortOrder) []Entry {
	result := make([]Entry, len(entries))
	copy(result, entries)

	switch order {
	case SortByStatus:
		sort.SliceStable(result, func(i, j int) bool {
			ri := statusRank[result[i].Status]
			rj := statusRank[result[j].Status]
			if ri != rj {
				return ri < rj
			}
			return result[i].Key < result[j].Key
		})
	default: // SortByKey
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Key < result[j].Key
		})
	}

	return result
}

// ParseSortOrder converts a string to a SortOrder value.
// Returns SortByKey and false if the value is unrecognized.
func ParseSortOrder(s string) (SortOrder, bool) {
	switch s {
	case "key":
		return SortByKey, true
	case "status":
		return SortByStatus, true
	default:
		return SortByKey, false
	}
}

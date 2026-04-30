package diff

// MergeResult holds the result of merging entries from a base env
// with suggested fixes derived from a reference env.
type MergeResult struct {
	Key      string
	OldValue string
	NewValue string
	Action   string // "add", "update", "keep"
}

// Merge produces a list of MergeResults describing what changes would be
// needed to bring base in line with reference, based on diff entries.
// Only Missing and Mismatch entries are acted upon; Match entries are kept.
func Merge(entries []Entry, base, reference map[string]string) []MergeResult {
	results := make([]MergeResult, 0, len(entries))

	for _, e := range entries {
		switch e.Status {
		case StatusMissing:
			results = append(results, MergeResult{
				Key:      e.Key,
				OldValue: "",
				NewValue: reference[e.Key],
				Action:   "add",
			})
		case StatusMismatch:
			results = append(results, MergeResult{
				Key:      e.Key,
				OldValue: base[e.Key],
				NewValue: reference[e.Key],
				Action:   "update",
			})
		case StatusMatch:
			results = append(results, MergeResult{
				Key:      e.Key,
				OldValue: base[e.Key],
				NewValue: base[e.Key],
				Action:   "keep",
			})
		}
	}

	return results
}

// ApplyMerge converts a slice of MergeResults into a flat key/value map
// representing the merged environment.
func ApplyMerge(results []MergeResult) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		if r.Action == "add" || r.Action == "update" || r.Action == "keep" {
			out[r.Key] = r.NewValue
		}
	}
	return out
}

package diff

// KeyStatus represents the comparison result for a single key.
type KeyStatus string

const (
	StatusMissing   KeyStatus = "missing"   // key exists in base but not in target
	StatusExtra     KeyStatus = "extra"     // key exists in target but not in base
	StatusMismatch  KeyStatus = "mismatch"  // key exists in both but values differ
	StatusMatch     KeyStatus = "match"     // key exists in both with equal values
)

// Result holds the comparison result for a single key.
type Result struct {
	Key        string
	Status     KeyStatus
	BaseValue  string
	TargetValue string
}

// Compare compares two parsed env maps (base vs target) and returns a slice
// of Result entries for every key found in either map.
func Compare(base, target map[string]string) []Result {
	results := make([]Result, 0)

	// Check keys in base
	for key, baseVal := range base {
		if targetVal, ok := target[key]; !ok {
			results = append(results, Result{
				Key:       key,
				Status:    StatusMissing,
				BaseValue: baseVal,
			})
		} else if baseVal != targetVal {
			results = append(results, Result{
				Key:         key,
				Status:      StatusMismatch,
				BaseValue:   baseVal,
				TargetValue: targetVal,
			})
		} else {
			results = append(results, Result{
				Key:         key,
				Status:      StatusMatch,
				BaseValue:   baseVal,
				TargetValue: targetVal,
			})
		}
	}

	// Check keys only in target
	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			results = append(results, Result{
				Key:         key,
				Status:      StatusExtra,
				TargetValue: targetVal,
			})
		}
	}

	return results
}

// FilterByStatus returns only the results matching one of the given statuses.
func FilterByStatus(results []Result, statuses ...KeyStatus) []Result {
	set := make(map[KeyStatus]struct{}, len(statuses))
	for _, s := range statuses {
		set[s] = struct{}{}
	}
	filtered := make([]Result, 0)
	for _, r := range results {
		if _, ok := set[r.Status]; ok {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

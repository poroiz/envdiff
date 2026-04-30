package diff

import "fmt"

// StatusFilter holds configuration for filtering diff entries.
type StatusFilter struct {
	IncludeMissing  bool
	IncludeExtra    bool
	IncludeMismatch bool
	IncludeMatch    bool
}

// DefaultFilter returns a StatusFilter that includes all problem statuses
// (missing, extra, mismatch) but excludes exact matches.
func DefaultFilter() StatusFilter {
	return StatusFilter{
		IncludeMissing:  true,
		IncludeExtra:    true,
		IncludeMismatch: true,
		IncludeMatch:    false,
	}
}

// AllFilter returns a StatusFilter that includes every status.
func AllFilter() StatusFilter {
	return StatusFilter{
		IncludeMissing:  true,
		IncludeExtra:    true,
		IncludeMismatch: true,
		IncludeMatch:    true,
	}
}

// Apply returns only the entries that match the StatusFilter settings.
func (f StatusFilter) Apply(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		switch e.Status {
		case StatusMissing:
			if f.IncludeMissing {
				out = append(out, e)
			}
		case StatusExtra:
			if f.IncludeExtra {
				out = append(out, e)
			}
		case StatusMismatch:
			if f.IncludeMismatch {
				out = append(out, e)
			}
		case StatusMatch:
			if f.IncludeMatch {
				out = append(out, e)
			}
		}
	}
	return out
}

// ParseStatusFilter builds a StatusFilter from a slice of status name strings.
// Accepted values: "missing", "extra", "mismatch", "match".
// An empty slice returns DefaultFilter.
func ParseStatusFilter(statuses []string) (StatusFilter, error) {
	if len(statuses) == 0 {
		return DefaultFilter(), nil
	}
	f := StatusFilter{}
	for _, s := range statuses {
		switch s {
		case "missing":
			f.IncludeMissing = true
		case "extra":
			f.IncludeExtra = true
		case "mismatch":
			f.IncludeMismatch = true
		case "match":
			f.IncludeMatch = true
		default:
			return StatusFilter{}, fmt.Errorf("unknown status filter: %q", s)
		}
	}
	return f, nil
}

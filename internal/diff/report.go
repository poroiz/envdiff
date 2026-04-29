package diff

import (
	"fmt"
	"io"
	"sort"
)

// PrintReport writes a human-readable diff report to w.
// It groups results by status and sorts keys alphabetically within each group.
func PrintReport(w io.Writer, results []Result, baseFile, targetFile string) {
	fmt.Fprintf(w, "Comparing %s (base) → %s (target)\n", baseFile, targetFile)
	fmt.Fprintln(w, "")

	groups := map[KeyStatus][]Result{
		StatusMissing:  FilterByStatus(results, StatusMissing),
		StatusExtra:    FilterByStatus(results, StatusExtra),
		StatusMismatch: FilterByStatus(results, StatusMismatch),
	}

	order := []KeyStatus{StatusMissing, StatusExtra, StatusMismatch}
	labels := map[KeyStatus]string{
		StatusMissing:  "MISSING in target",
		StatusExtra:    "EXTRA in target",
		StatusMismatch: "VALUE MISMATCH",
	}

	anyIssue := false
	for _, status := range order {
		items := groups[status]
		if len(items) == 0 {
			continue
		}
		anyIssue = true
		sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })
		fmt.Fprintf(w, "[%s]\n", labels[status])
		for _, r := range items {
			switch r.Status {
			case StatusMissing:
				fmt.Fprintf(w, "  - %s (base: %q)\n", r.Key, r.BaseValue)
			case StatusExtra:
				fmt.Fprintf(w, "  + %s (target: %q)\n", r.Key, r.TargetValue)
			case StatusMismatch:
				fmt.Fprintf(w, "  ~ %s\n      base:   %q\n      target: %q\n", r.Key, r.BaseValue, r.TargetValue)
			}
		}
		fmt.Fprintln(w, "")
	}

	if !anyIssue {
		fmt.Fprintln(w, "✓ No differences found.")
	}

	match := FilterByStatus(results, StatusMatch)
	fmt.Fprintf(w, "Summary: %d matching, %d missing, %d extra, %d mismatched\n",
		len(match),
		len(groups[StatusMissing]),
		len(groups[StatusExtra]),
		len(groups[StatusMismatch]),
	)
}

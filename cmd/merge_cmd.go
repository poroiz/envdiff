package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func newMergeCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "merge [base] [reference]",
		Short: "Show changes needed to bring base env in line with reference",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "preview changes without writing (default true)")
	return cmd
}

func runMerge(basePath, refPath string, dryRun bool) error {
	base, err := parser.ParseFile(basePath)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}

	ref, err := parser.ParseFile(refPath)
	if err != nil {
		return fmt.Errorf("reading reference file: %w", err)
	}

	entries := diff.Compare(base, ref)
	results := diff.Merge(entries, base, ref)

	// Sort results by key for deterministic output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	w := os.Stdout
	fmt.Fprintf(w, "Merge preview: %s <- %s\n\n", basePath, refPath)

	for _, r := range results {
		switch r.Action {
		case "add":
			fmt.Fprintf(w, "  [ADD]    %s=%s\n", r.Key, r.NewValue)
		case "update":
			fmt.Fprintf(w, "  [UPDATE] %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
		case "keep":
			fmt.Fprintf(w, "  [KEEP]   %s=%s\n", r.Key, r.NewValue)
		}
	}

	if dryRun {
		fmt.Fprintln(w, "\n(dry-run mode: no files were modified)")
	}

	return nil
}

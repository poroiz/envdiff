package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/parser"
)

var formatFlag string

// newCompareCmd builds the compare subcommand which supports --format.
func newCompareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <base> <other>",
		Short: "Compare two .env files and report differences",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	cmd.Flags().StringVarP(&formatFlag, "format", "f", "text",
		"Output format: text, json, or csv")
	return cmd
}

func runCompare(cmd *cobra.Command, args []string) error {
	baseFile := args[0]
	otherFile := args[1]

	fmt, err := output.ParseFormat(formatFlag)
	if err != nil {
		return err
	}

	baseMap, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt_errorf("reading base file: %w", err)
	}

	otherMap, err := parser.ParseFile(otherFile)
	if err != nil {
		return fmt_errorf("reading compare file: %w", err)
	}

	entries := diff.Compare(baseMap, otherMap)

	if err := output.Write(os.Stdout, entries, fmt); err != nil {
		return err
	}

	// Exit with code 1 if any non-match entries exist.
	for _, e := range entries {
		if e.Status != diff.StatusMatch {
			os.Exit(1)
		}
	}
	return nil
}

// fmt_errorf is a thin wrapper to avoid import collision with the output.Format type.
func fmt_errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

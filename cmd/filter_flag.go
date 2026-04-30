package cmd

import (
	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/output"
	"envdiff/internal/parser"
)

var filterStatuses []string

func newFilterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter <base> <compare>",
		Short: "Compare .env files and display only entries matching the given statuses",
		Args:  cobra.ExactArgs(2),
		RunE:  runFilter,
	}
	cmd.Flags().StringSliceVar(
		&filterStatuses, "status", []string{},
		"Comma-separated list of statuses to include: missing, extra, mismatch, match (default: missing,extra,mismatch)",
	)
	cmd.Flags().StringP("format", "f", "text", "Output format: text, json, csv")
	return cmd
}

func runFilter(cmd *cobra.Command, args []string) error {
	baseFile, compareFile := args[0], args[1]

	baseMap, err := parser.ParseFile(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}
	cmpMap, err := parser.ParseFile(compareFile)
	if err != nil {
		return fmt.Errorf("reading compare file: %w", err)
	}

	entries := diff.Compare(baseMap, cmpMap)

	f, err := diff.ParseStatusFilter(filterStatuses)
	if err != nil {
		return err
	}
	filtered := f.Apply(entries)

	fmtFlag, _ := cmd.Flags().GetString("format")
	fmt, err := output.ParseFormat(fmtFlag)
	if err != nil {
		return err
	}

	return output.Write(cmd.OutOrStdout(), filtered, fmt)
}

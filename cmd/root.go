package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	statusFilter string
	silent       bool
)

var rootCmd = &cobra.Command{
	Use:   "envdiff <base> <compare>",
	Short: "Compare .env files across environments",
	Long: `envdiff compares two .env files and reports missing,
extra, or mismatched keys between them.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseFile := args[0]
		cmpFile := args[1]

		baseEnv, err := parser.ParseFile(baseFile)
		if err != nil {
			return fmt.Errorf("failed to parse base file %q: %w", baseFile, err)
		}

		cmpEnv, err := parser.ParseFile(cmpFile)
		if err != nil {
			return fmt.Errorf("failed to parse compare file %q: %w", cmpFile, err)
		}

		results := diff.Compare(baseEnv, cmpEnv)

		if statusFilter != "" {
			results = diff.FilterByStatus(results, statusFilter)
		}

		if !silent {
			diff.PrintReport(results, baseFile, cmpFile)
		}

		// Exit with non-zero code if any differences found
		for _, r := range results {
			if r.Status != "match" {
				os.Exit(1)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&statusFilter, "filter", "f", "",
		"filter output by status: match, missing, extra, mismatch")
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false,
		"suppress output, only use exit code")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

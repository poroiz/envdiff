package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format defines the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// ParseFormat converts a string to a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be one of text, json, csv", s)
	}
}

// Write renders diff entries to w using the specified format.
func Write(w io.Writer, entries []diff.Entry, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, entries)
	case FormatCSV:
		return writeCSV(w, entries)
	default:
		return writeText(w, entries)
	}
}

func writeText(w io.Writer, entries []diff.Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "[%s] %s\n", e.Status, e.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []diff.Entry) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

func writeCSV(w io.Writer, entries []diff.Entry) error {
	_, err := fmt.Fprintln(w, "key,status,base_value,compare_value")
	if err != nil {
		return err
	}
	for _, e := range entries {
		_, err := fmt.Fprintf(w, "%s,%s,%s,%s\n",
			csvEscape(e.Key),
			string(e.Status),
			csvEscape(e.BaseValue),
			csvEscape(e.CompareValue),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

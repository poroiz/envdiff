package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/output"
)

var sampleEntries = []diff.Entry{
	{Key: "DB_HOST", Status: diff.StatusMatch, BaseValue: "localhost", CompareValue: "localhost"},
	{Key: "DB_PASS", Status: diff.StatusMismatch, BaseValue: "secret", CompareValue: "other"},
	{Key: "API_KEY", Status: diff.StatusMissing, BaseValue: "abc123", CompareValue: ""},
	{Key: "NEW_VAR", Status: diff.StatusExtra, BaseValue: "", CompareValue: "xyz"},
}

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct{ in string; want output.Format }{
		{"text", output.FormatText},
		{"json", output.FormatJSON},
		{"csv", output.FormatCSV},
		{"", output.FormatText},
		{"JSON", output.FormatJSON},
	}
	for _, c := range cases {
		got, err := output.ParseFormat(c.in)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", c.in, err)
		}
		if got != c.want {
			t.Errorf("ParseFormat(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestWrite_Text(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleEntries, output.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[MISMATCH] DB_PASS") {
		t.Errorf("expected mismatch line, got:\n%s", out)
	}
}

func TestWrite_TextEmpty(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, []diff.Entry{}, output.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleEntries, output.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"Key\"") && !strings.Contains(out, "\"key\"") {
		t.Errorf("expected JSON output with keys, got:\n%s", out)
	}
}

func TestWrite_CSV(t *testing.T) {
	var buf bytes.Buffer
	if err := output.Write(&buf, sampleEntries, output.FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "key,status,base_value,compare_value" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if len(lines) != 5 {
		t.Errorf("expected 5 lines (header + 4 entries), got %d", len(lines))
	}
}

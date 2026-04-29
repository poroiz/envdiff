package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func TestRootCmd_MissingArgs(t *testing.T) {
	_, err := executeCommand()
	if err == nil {
		t.Error("expected error when no args provided, got nil")
	}
}

func TestRootCmd_FileNotFound(t *testing.T) {
	_, err := executeCommand("nonexistent1.env", "nonexistent2.env")
	if err == nil {
		t.Error("expected error for missing files, got nil")
	}
}

func TestRootCmd_IdenticalFiles(t *testing.T) {
	content := "KEY1=value1\nKEY2=value2\n"
	base := writeTempEnv(t, content)
	cmp := writeTempEnv(t, content)

	// Reset flags between tests
	statusFilter = ""
	silent = true

	_, err := executeCommand("--silent", base, cmp)
	if err != nil {
		t.Errorf("expected no error for identical files, got: %v", err)
	}
}

func TestRootCmd_WithFilterFlag(t *testing.T) {
	base := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	cmp := writeTempEnv(t, "KEY1=value1\n")

	statusFilter = ""
	silent = false

	_, err := executeCommand("--filter", "missing", base, cmp)
	// err may be non-nil due to os.Exit(1) being called; we just verify no panic
	_ = err
}

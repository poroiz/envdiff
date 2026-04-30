package diff

import (
	"testing"
)

func makeMergeBase() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DEBUG":    "false",
		"PORT":     "8080",
	}
}

func makeMergeRef() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"DEBUG":    "true",
		"LOG_LEVEL": "info",
	}
}

func TestMerge_AddsMissingKeys(t *testing.T) {
	base := makeMergeBase()
	ref := makeMergeRef()
	entries := Compare(base, ref)

	results := Merge(entries, base, ref)

	var added []MergeResult
	for _, r := range results {
		if r.Action == "add" {
			added = append(added, r)
		}
	}
	if len(added) != 1 || added[0].Key != "LOG_LEVEL" {
		t.Errorf("expected 1 add for LOG_LEVEL, got %+v", added)
	}
}

func TestMerge_UpdatesMismatchedKeys(t *testing.T) {
	base := makeMergeBase()
	ref := makeMergeRef()
	entries := Compare(base, ref)

	results := Merge(entries, base, ref)

	for _, r := range results {
		if r.Key == "DEBUG" {
			if r.Action != "update" {
				t.Errorf("expected action=update for DEBUG, got %s", r.Action)
			}
			if r.NewValue != "true" {
				t.Errorf("expected NewValue=true, got %s", r.NewValue)
			}
			return
		}
	}
	t.Error("DEBUG key not found in merge results")
}

func TestMerge_KeepsMatchingKeys(t *testing.T) {
	base := makeMergeBase()
	ref := makeMergeRef()
	entries := Compare(base, ref)

	results := Merge(entries, base, ref)

	for _, r := range results {
		if r.Key == "APP_NAME" && r.Action != "keep" {
			t.Errorf("expected action=keep for APP_NAME, got %s", r.Action)
		}
	}
}

func TestApplyMerge_ProducesCorrectMap(t *testing.T) {
	base := makeMergeBase()
	ref := makeMergeRef()
	entries := Compare(base, ref)
	results := Merge(entries, base, ref)

	merged := ApplyMerge(results)

	if merged["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %s", merged["LOG_LEVEL"])
	}
	if merged["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true after merge, got %s", merged["DEBUG"])
	}
	if merged["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %s", merged["APP_NAME"])
	}
}

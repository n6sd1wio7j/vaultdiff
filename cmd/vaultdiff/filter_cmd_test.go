package main

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/diff"
)

func TestFilterFlags_Default(t *testing.T) {
	ff := FilterFlags{}
	if ff.OnlyChanged {
		t.Error("OnlyChanged should default to false")
	}
	if ff.KeyPrefix != "" {
		t.Error("KeyPrefix should default to empty string")
	}
	if len(ff.ExcludeKeys) != 0 {
		t.Error("ExcludeKeys should default to empty")
	}
}

func TestFilterFlags_OnlyChangedFiltersUnchanged(t *testing.T) {
	entries := []diff.Entry{
		{Key: "x", Status: diff.StatusUnchanged},
		{Key: "y", Status: diff.StatusAdded},
	}
	result := diff.Filter(entries, diff.FilterOptions{OnlyChanged: true})
	if len(result) != 1 || result[0].Key != "y" {
		t.Errorf("expected only 'y', got %+v", result)
	}
}

func TestFilterFlags_PrefixIsolatesKeys(t *testing.T) {
	entries := []diff.Entry{
		{Key: "aws_key", Status: diff.StatusModified},
		{Key: "aws_secret", Status: diff.StatusModified},
		{Key: "gcp_token", Status: diff.StatusAdded},
	}
	result := diff.Filter(entries, diff.FilterOptions{KeyPrefix: "aws_"})
	if len(result) != 2 {
		t.Fatalf("expected 2 aws_ entries, got %d", len(result))
	}
}

func TestFilterFlags_ExcludeRemovesSensitiveKeys(t *testing.T) {
	entries := []diff.Entry{
		{Key: "password", Status: diff.StatusModified},
		{Key: "username", Status: diff.StatusUnchanged},
	}
	result := diff.Filter(entries, diff.FilterOptions{ExcludeKeys: []string{"password"}})
	if len(result) != 1 || result[0].Key != "username" {
		t.Errorf("expected only 'username', got %+v", result)
	}
}

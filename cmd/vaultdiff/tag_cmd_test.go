package main

import (
	"testing"
)

func TestTagFlags_MissingPath(t *testing.T) {
	err := RunTag([]string{"--addr", "http://vault", "--token", "tok"})
	if err == nil || err.Error() != "--path is required" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestTagFlags_MissingAddr(t *testing.T) {
	err := RunTag([]string{"--path", "secret/app", "--token", "tok"})
	if err == nil || err.Error() != "--addr is required" {
		t.Errorf("expected addr error, got %v", err)
	}
}

func TestTagFlags_MissingToken(t *testing.T) {
	err := RunTag([]string{"--path", "secret/app", "--addr", "http://vault"})
	if err == nil || err.Error() != "--token is required" {
		t.Errorf("expected token error, got %v", err)
	}
}

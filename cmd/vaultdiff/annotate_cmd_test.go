package main

import (
	"testing"
)

func TestAnnotateFlags_MissingPath(t *testing.T) {
	err := RunAnnotate([]string{"--addr", "http://vault", "--token", "tok"})
	if err == nil || err.Error() != "--path is required" {
		t.Errorf("expected path error, got %v", err)
	}
}

func TestAnnotateFlags_MissingAddr(t *testing.T) {
	err := RunAnnotate([]string{"--path", "secret/foo", "--token", "tok"})
	if err == nil || err.Error() != "--addr is required" {
		t.Errorf("expected addr error, got %v", err)
	}
}

func TestAnnotateFlags_MissingToken(t *testing.T) {
	err := RunAnnotate([]string{"--path", "secret/foo", "--addr", "http://vault"})
	if err == nil || err.Error() != "--token is required" {
		t.Errorf("expected token error, got %v", err)
	}
}

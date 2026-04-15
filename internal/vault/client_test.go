package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_DefaultMount(t *testing.T) {
	client, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.MountPath != "secret" {
		t.Errorf("expected default mount 'secret', got %q", client.MountPath)
	}
}

func TestReadSecretVersion_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Address:   server.URL,
		Token:     "test-token",
		MountPath: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	_, err = client.ReadSecretVersion("myapp/config", 0)
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}

func TestReadSecretVersion_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"API_KEY":"abc123","DB_PASS":"secret"}}}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Address:   server.URL,
		Token:     "test-token",
		MountPath: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	data, err := client.ReadSecretVersion("myapp/config", 1)
	if err != nil {
		t.Fatalf("unexpected error reading secret: %v", err)
	}
	if data["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %v", data["API_KEY"])
	}
}

package vault

import (
	"context"
	"errors"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestListVersions_NoMetadata(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return nil, nil
		},
	}
	c := &Client{logical: mock, mount: "secret"}
	_, err := c.ListVersions(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error for nil secret, got nil")
	}
}

func TestListVersions_ReadError(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return nil, errors.New("connection refused")
		},
	}
	c := &Client{logical: mock, mount: "secret"}
	_, err := c.ListVersions(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error on read failure, got nil")
	}
}

func TestListVersions_Success(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"versions": map[string]interface{}{
						"1": map[string]interface{}{
							"created_time":  "2024-01-01T00:00:00Z",
							"deletion_time": "",
							"destroyed":     false,
						},
						"2": map[string]interface{}{
							"created_time":  "2024-02-01T00:00:00Z",
							"deletion_time": "",
							"destroyed":     false,
						},
					},
				},
			}, nil
		},
	}
	c := &Client{logical: mock, mount: "secret"}
	metas, err := c.ListVersions(context.Background(), "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(metas))
	}
}

func TestListVersions_MissingVersionsKey(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"other_key": "value",
				},
			}, nil
		},
	}
	c := &Client{logical: mock, mount: "secret"}
	_, err := c.ListVersions(context.Background(), "myapp/config")
	if err == nil {
		t.Fatal("expected error for missing versions key, got nil")
	}
}

package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with vaultdiff-specific helpers.
type Client struct {
	api     *vaultapi.Client
	MountPath string
}

// Config holds connection settings for a Vault instance.
type Config struct {
	Address   string
	Token     string
	MountPath string
}

// NewClient creates and configures a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address is required (set VAULT_ADDR or pass --address)")
	}
	vcfg.Address = address

	client, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set VAULT_TOKEN or pass --token)")
	}
	client.SetToken(token)

	mount := cfg.MountPath
	if mount == "" {
		mount = "secret"
	}

	return &Client{api: client, MountPath: mount}, nil
}

// ReadSecretVersion reads a specific version of a KV v2 secret.
// Pass version 0 to read the latest version.
func (c *Client) ReadSecretVersion(secretPath string, version int) (map[string]interface{}, error) {
	params := map[string][]string{}
	if version > 0 {
		params["version"] = []string{fmt.Sprintf("%d", version)}
	}

	kvPath := fmt.Sprintf("%s/data/%s", c.MountPath, secretPath)
	secret, err := c.api.Logical().ReadWithData(kvPath, params)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret %q: %w", secretPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret %q not found", secretPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for secret %q", secretPath)
	}
	return data, nil
}

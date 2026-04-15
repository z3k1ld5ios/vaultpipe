package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	vc *vaultapi.Client
}

// Config holds the configuration needed to connect to Vault.
type Config struct {
	Address string
	Token   string
	RoleID  string
	SecretID string
}

// NewClient creates a new Vault client from the provided config.
// If Address or Token are empty, it falls back to environment variables
// VAULT_ADDR and VAULT_TOKEN respectively.
func NewClient(cfg Config) (*Client, error) {
	vcCfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr == "" {
		return nil, fmt.Errorf("vault address is required: set --vault-addr or VAULT_ADDR")
	}
	vcCfg.Address = addr

	vc, err := vaultapi.NewClient(vcCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token != "" {
		vc.SetToken(token)
	}

	return &Client{vc: vc}, nil
}

// ReadSecretKV2 reads a KV v2 secret at the given mount and path,
// returning the key/value data map.
func (c *Client) ReadSecretKV2(mount, secretPath string) (map[string]string, error) {
	fullPath := fmt.Sprintf("%s/data/%s", mount, secretPath)
	secret, err := c.vc.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", fullPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at path %q", fullPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("secret key %q has non-string value", k)
		}
		result[k] = str
	}
	return result, nil
}

package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretMap is a flat map of key-value secret pairs.
type SecretMap map[string]string

// ReadSecretKV2 reads a KV v2 secret from Vault at the given mount and path,
// returning a flat map of the secret's data fields.
func (c *Client) ReadSecretKV2(ctx context.Context, mount, path string) (SecretMap, error) {
	if mount == "" {
		mount = "secret"
	}
	path = strings.TrimPrefix(path, "/")
	secretPath := fmt.Sprintf("%s/data/%s", mount, path)

	secret, err := c.logical.ReadWithContext(ctx, secretPath)
	if err != nil {
		return nil, fmt.Errorf("vault read %q: %w", secretPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault secret not found at %q", secretPath)
	}

	data, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("vault secret at %q missing 'data' field", secretPath)
	}

	raw, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault secret data at %q has unexpected type", secretPath)
	}

	result := make(SecretMap, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			result[k] = val
		default:
			result[k] = fmt.Sprintf("%v", val)
		}
	}
	return result, nil
}

// ReadMultiple reads secrets from multiple paths and merges them into a single
// SecretMap. Later paths take precedence over earlier ones on key conflicts.
func (c *Client) ReadMultiple(ctx context.Context, mount string, paths []string) (SecretMap, error) {
	merged := make(SecretMap)
	for _, path := range paths {
		sm, err := c.ReadSecretKV2(ctx, mount, path)
		if err != nil {
			return nil, fmt.Errorf("reading path %q: %w", path, err)
		}
		for k, v := range sm {
			merged[k] = v
		}
	}
	return merged, nil
}

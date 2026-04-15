// Package token manages Vault token lifecycle including retrieval and validation.
package token

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Source defines where a Vault token originates.
type Source string

const (
	SourceEnv    Source = "env"
	SourceFile   Source = "file"
	SourceStatic Source = "static"
)

// Token holds a resolved Vault token and its origin.
type Token struct {
	Value  string
	Source Source
}

// Provider resolves a Vault token from configured sources.
type Provider struct {
	staticToken string
	tokenFile   string
	envVar      string
}

// NewProvider creates a Provider. envVar defaults to VAULT_TOKEN if empty.
func NewProvider(staticToken, tokenFile, envVar string) *Provider {
	if envVar == "" {
		envVar = "VAULT_TOKEN"
	}
	return &Provider{
		staticToken: staticToken,
		tokenFile:   tokenFile,
		envVar:      envVar,
	}
}

// Resolve attempts to obtain a token, checking static > env > file order.
func (p *Provider) Resolve(_ context.Context) (*Token, error) {
	if p.staticToken != "" {
		return &Token{Value: p.staticToken, Source: SourceStatic}, nil
	}

	if val := os.Getenv(p.envVar); val != "" {
		return &Token{Value: val, Source: SourceEnv}, nil
	}

	if p.tokenFile != "" {
		data, err := os.ReadFile(p.tokenFile)
		if err != nil {
			return nil, fmt.Errorf("token: reading file %q: %w", p.tokenFile, err)
		}
		v := strings.TrimSpace(string(data))
		if v == "" {
			return nil, errors.New("token: file is empty")
		}
		return &Token{Value: v, Source: SourceFile}, nil
	}

	return nil, errors.New("token: no token source configured (static, env, or file)")
}

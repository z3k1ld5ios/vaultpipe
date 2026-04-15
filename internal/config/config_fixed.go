package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SecretMapping defines a single Vault secret path and how its keys
// should be mapped into environment variable names.
type SecretMapping struct {
	Path   string            `yaml:"path"`
	Mount  string            `yaml:"mount"`
	EnvMap map[string]string `yaml:"env_map"`
}

// Config holds the top-level vaultpipe configuration.
type Config struct {
	VaultAddress string          `yaml:"vault_address"`
	Secrets      []SecretMapping `yaml:"secrets"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded config.
func (c *Config) validate() error {
	if len(c.Secrets) == 0 {
		return errors.New("config must define at least one secret mapping")
	}
	for i, s := range c.Secrets {
		if s.Path == "" {
			return fmt.Errorf("secret mapping[%d]: path must not be empty", i)
		}
		if s.Mount == "" {
			c.Secrets[i].Mount = "secret"
		}
	}
	return nil
}

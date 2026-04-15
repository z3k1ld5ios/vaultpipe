# vaultpipe

> A lightweight CLI for piping secrets from HashiCorp Vault into process environments without writing to disk.

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

---

## Usage

Inject secrets from a Vault path directly into a process environment:

```bash
vaultpipe run --path secret/data/myapp -- ./myapp
```

The secrets stored at the given Vault path are exported as environment variables to the child process. Nothing is written to disk.

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--prefix` | Env var name prefix | *(none)* |

**Example with prefix:**

```bash
vaultpipe run --path secret/data/db --prefix DB_ -- ./server
```

This would expose `username` as `DB_USERNAME`, `password` as `DB_PASSWORD`, and so on.

---

## Authentication

`vaultpipe` respects standard Vault environment variables (`VAULT_ADDR`, `VAULT_TOKEN`, `VAULT_NAMESPACE`) and uses the local Vault agent token if available.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

---

## License

[MIT](LICENSE)
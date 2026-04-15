# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments.

---

## Installation

```bash
go install github.com/youruser/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/vaultdiff.git
cd vaultdiff
go build -o vaultdiff .
```

---

## Usage

```bash
# Diff two versions of a secret
vaultdiff --path secret/myapp/config --v1 3 --v2 4

# Compare secrets across environments
vaultdiff --src secret/staging/myapp --dst secret/production/myapp

# Output diff in JSON format
vaultdiff --path secret/myapp/config --v1 2 --v2 5 --format json
```

**Example output:**

```diff
~ DB_HOST: "db-staging.internal" → "db-prod.internal"
+ NEW_FEATURE_FLAG: "true"
- DEPRECATED_KEY: "old-value"
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | — |
| `--v1` | First version to compare | latest-1 |
| `--v2` | Second version to compare | latest |
| `--src` | Source path for cross-env diff | — |
| `--dst` | Destination path for cross-env diff | — |
| `--format` | Output format (`text`, `json`) | `text` |

### Authentication

`vaultdiff` respects standard Vault environment variables:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"
```

---

## License

[MIT](LICENSE)
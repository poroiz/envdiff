# envdiff

A CLI tool that compares `.env` files across environments and flags missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-env> <compare-env> [additional-envs...]
```

### Example

```bash
envdiff .env.example .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED keys (present in both, different values):
  - APP_ENV  (.env.example: "development" | .env.production: "production")

✔ All other keys match.
```

### Flags

| Flag | Description |
|------|-------------|
| `--strict` | Exit with non-zero status if any differences are found |
| `--ignore-values` | Only check for missing keys, skip value comparison |
| `--quiet` | Suppress output, use exit code only |

---

## Use Cases

- Validate `.env` files in CI/CD pipelines before deployment
- Onboard new developers by comparing their local `.env` against `.env.example`
- Audit environment parity across staging and production

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE)
# Environment Variables

RevenueCat CLI can be configured entirely through environment variables. This is especially useful for CI/CD pipelines, Docker containers, and automated workflows.

## Available Variables

| Variable | Description | Example | Default |
|---|---|---|---|
| `RC_API_KEY` | Secret API key (overrides active profile) | `sk_xxxxx` | -- |
| `RC_PROJECT` | Default project ID | `proj_xxxxx` | -- |
| `RC_PROFILE` | Active auth profile name | `production` | First created profile |
| `RC_OUTPUT` | Default output format | `table` | `json` |
| `RC_DEBUG` | Enable debug logging (shows HTTP requests/responses) | `true` | `false` |
| `RC_TIMEOUT` | HTTP request timeout duration | `30s` | `10s` |

---

## Precedence

Environment variables sit in the middle of the configuration hierarchy:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Local `.rc.yaml` file**
4. **Active auth profile**
5. **Built-in defaults** (lowest priority)

### Example

```bash
export RC_OUTPUT=yaml

# Uses YAML (from env var)
rc apps list

# Uses table (flag overrides env var)
rc apps list --output table
```

---

## Usage Examples

### CI/CD Pipeline

```bash
export RC_API_KEY="sk_xxxxx"
export RC_PROJECT="proj_xxxxx"
export RC_OUTPUT="json"

# No auth profile needed -- env vars are sufficient
rc products list | jq '.items | length'
rc entitlements list --all --output csv > entitlements.csv
```

### Docker

```dockerfile
ENV RC_API_KEY=sk_xxxxx
ENV RC_PROJECT=proj_xxxxx
ENV RC_TIMEOUT=30s
```

### GitHub Actions

```yaml
steps:
  - name: Check subscription metrics
    env:
      RC_API_KEY: ${{ secrets.REVENUECAT_API_KEY }}
      RC_PROJECT: proj_xxxxx
    run: |
      rc metrics overview --output json
```

### Shell Profile

Add to `~/.bashrc` or `~/.zshrc` for persistent defaults:

```bash
export RC_PROFILE="production"
export RC_OUTPUT="table"
```

---

## Debugging

Enable debug mode to see full HTTP request and response details:

```bash
export RC_DEBUG=true
rc apps list
```

This outputs request URLs, headers, response status codes, and timing information. Useful for troubleshooting API errors.

!!! warning
    Debug mode prints API responses to stderr, which may include sensitive data. Do not enable it in production logs.

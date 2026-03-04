# Configuration

RevenueCat CLI supports multiple authentication profiles, project-level defaults, and environment variable overrides. This page covers all configuration options.

## Authentication Profiles

Profiles let you store multiple API keys and switch between them. This is useful when managing staging and production environments, or multiple RevenueCat accounts.

### Create a Profile

```bash
rc auth login --api-key sk_xxxxx --name production
```

You can optionally bind a default project to a profile:

```bash
rc auth login --api-key sk_xxxxx --name production --default-project proj_xxxxx
```

### Switch Profiles

```bash
rc auth switch --name staging
```

### List All Profiles

```bash
rc auth list
```

### View Current Profile

```bash
rc auth current
```

### Delete a Profile

```bash
rc auth delete --name old-profile --confirm
```

!!! warning
    Deleting a profile removes its stored API key permanently. This action cannot be undone.

---

## Project Configuration

### Initialize a Project

Bind a default project to the current directory:

```bash
rc init --project proj_xxxxx
```

This creates a `.rc.yaml` file in the current directory.

### `.rc.yaml` File Format

```yaml
project: proj_xxxxx
```

When a `.rc.yaml` file is present, all commands automatically use that project ID. You can override it per-command with the `--project` flag.

!!! tip
    Commit `.rc.yaml` to your repository so your team shares the same project configuration.

---

## Environment Variables

Environment variables override profile and file-based configuration. They are especially useful in CI/CD pipelines.

| Variable | Description | Example |
|---|---|---|
| `RC_API_KEY` | API secret key (overrides active profile) | `sk_xxxxx` |
| `RC_PROJECT` | Default project ID | `proj_xxxxx` |
| `RC_PROFILE` | Active profile name | `production` |
| `RC_OUTPUT` | Default output format | `json`, `table`, `csv` |
| `RC_DEBUG` | Enable debug logging | `true` |
| `RC_TIMEOUT` | HTTP request timeout | `30s` |

### Precedence Order

Configuration is resolved in this order (highest priority first):

1. Command-line flags (`--project`, `--output`, etc.)
2. Environment variables (`RC_PROJECT`, `RC_OUTPUT`, etc.)
3. Local `.rc.yaml` file
4. Active auth profile defaults
5. Built-in defaults

!!! note
    Flags always win. If you pass `--project proj_yyy` on the command line, it overrides everything else.

## Next Steps

Head to the [Quick Start](quickstart.md) to run your first commands.

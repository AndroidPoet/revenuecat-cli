# Projects

Manage RevenueCat projects. A project is the top-level container for your apps, products, and entitlements.

## Available Commands

| Command | Description |
|---|---|
| `rc projects list` | List all projects accessible with your API key |
| `rc projects create` | Create a new project |

---

## `rc projects list`

List all projects associated with the current API key.

```bash
rc projects list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format (`json`, `pretty`, `table`, `csv`, `tsv`, `yaml`, `minimal`) | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

### Example

```bash
rc projects list --output table
```

```
ID            NAME
proj_xxxxx    Production
proj_yyyyy    Staging
```

---

## `rc projects create`

Create a new RevenueCat project.

```bash
rc projects create --name "My New Project"
```

### Required Flags

| Flag | Description |
|---|---|
| `--name` | Name for the new project |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

### Example

```bash
rc projects create --name "Staging Environment" --output pretty
```

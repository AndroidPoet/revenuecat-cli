# Output Formats

RevenueCat CLI supports 7 output formats. Use the `--output` flag (or the `RC_OUTPUT` environment variable) to switch between them.

```bash
rc <command> --output <format>
```

---

## Available Formats

| Format | Flag Value | Best For |
|---|---|---|
| JSON | `json` | Piping to `jq`, API integration, scripting |
| Pretty JSON | `pretty` | Human-readable inspection |
| Table | `table` | Quick visual overview in the terminal |
| CSV | `csv` | Spreadsheet import, data analysis |
| TSV | `tsv` | Tab-separated processing, clipboard pasting |
| YAML | `yaml` | Configuration files, readability |
| Minimal | `minimal` | Extracting single values for shell scripts |

---

## Format Examples

All examples use `rc apps list` with two apps in the project.

### JSON (default)

```bash
rc apps list --output json
```

```json
{
  "items": [
    {
      "id": "app_xxxxx",
      "name": "My iOS App",
      "type": "app_store"
    },
    {
      "id": "app_yyyyy",
      "name": "My Android App",
      "type": "play_store"
    }
  ],
  "next_page": null
}
```

### Pretty JSON

```bash
rc apps list --output pretty
```

Identical structure to JSON but with syntax highlighting and indentation optimized for terminal readability.

### Table

```bash
rc apps list --output table
```

```
ID          NAME              TYPE
app_xxxxx   My iOS App        app_store
app_yyyyy   My Android App    play_store
```

### CSV

```bash
rc apps list --output csv
```

```
id,name,type
app_xxxxx,My iOS App,app_store
app_yyyyy,My Android App,play_store
```

### TSV

```bash
rc apps list --output tsv
```

```
id	name	type
app_xxxxx	My iOS App	app_store
app_yyyyy	My Android App	play_store
```

### YAML

```bash
rc apps list --output yaml
```

```yaml
items:
  - id: app_xxxxx
    name: My iOS App
    type: app_store
  - id: app_yyyyy
    name: My Android App
    type: play_store
next_page: null
```

### Minimal

```bash
rc apps list --output minimal
```

```
app_xxxxx My iOS App app_store
app_yyyyy My Android App play_store
```

!!! tip
    Minimal output is ideal for shell scripts. Combine with `awk` to extract specific fields:

    ```bash
    rc apps list --output minimal | awk '{print $1}'
    ```

---

## Setting a Default Format

### Per-session

```bash
export RC_OUTPUT=table
rc apps list          # uses table format
rc products list      # also uses table format
```

### Per-profile

```bash
rc auth login --api-key sk_xxxxx --name dev
# Then set RC_OUTPUT in your shell profile for that context
```

### Per-command

```bash
rc apps list --output csv
```

The `--output` flag always takes precedence over the environment variable.

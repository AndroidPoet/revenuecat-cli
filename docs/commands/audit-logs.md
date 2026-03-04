# Audit Logs

View the audit log for your RevenueCat project. Audit logs track all administrative actions -- configuration changes, API key usage, team member activity, and more.

## Available Commands

| Command | Description |
|---|---|
| `rc audit-logs list` | List audit log entries |

---

## `rc audit-logs list`

Retrieve audit log entries with support for cursor-based pagination.

```bash
rc audit-logs list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results per page | `20` |
| `--starting-after` | Cursor for pagination (use the last entry's ID) | -- |
| `--all` | Fetch all pages automatically | `false` |

### Example

```bash
rc audit-logs list --limit 10 --output table
```

```
ID            ACTION              ACTOR           TIMESTAMP
aud_xxxxx     entitlement.create  user@example    2025-01-15T10:30:00Z
aud_yyyyy     offering.update     api_key_xxxx    2025-01-15T09:15:00Z
aud_zzzzz     app.delete          user@example    2025-01-14T18:45:00Z
```

### Paginating Through Results

Audit logs can be large. Use cursor-based pagination to walk through results page by page:

```bash
# First page
rc audit-logs list --limit 50

# Next page (use the last ID from the previous response)
rc audit-logs list --limit 50 --starting-after aud_last_id

# Or fetch everything at once
rc audit-logs list --all --output csv > audit-log-export.csv
```

!!! tip
    Export the full audit log to CSV for compliance reporting: `rc audit-logs list --all --output csv > audit.csv`

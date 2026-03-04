# Quick Start

Get up and running with RevenueCat CLI in under a minute.

## Step 1: Authenticate

```bash
rc auth login --api-key sk_xxxxx --name myproject
```

This stores your API key in a named profile. You can create multiple profiles for different environments.

## Step 2: Set Your Project

```bash
rc init --project proj_xxxxx
```

This writes a `.rc.yaml` file in the current directory so all subsequent commands know which project to target.

## Step 3: Verify

```bash
rc doctor
```

The doctor command validates your authentication, project configuration, and API connectivity. A successful check looks like:

```
Authentication:  OK
Project:         OK (proj_xxxxx)
API Connection:  OK (latency: 120ms)
```

---

## Your First Commands

### List Your Apps

```bash
rc apps list
```

### List Products in Table Format

```bash
rc products list --output table
```

### Get a Specific Customer

```bash
rc customers get --customer-id cust_xxxxx --output pretty
```

### View Metrics Overview

```bash
rc metrics overview --output table
```

---

## Output Format Examples

RevenueCat CLI supports 7 output formats. Use the `--output` flag to switch between them.

=== "JSON (default)"

    ```bash
    rc apps list --output json
    ```

    ```json
    {
      "items": [
        {
          "id": "app_xxxxx",
          "name": "My App",
          "type": "app_store"
        }
      ]
    }
    ```

=== "Table"

    ```bash
    rc apps list --output table
    ```

    ```
    ID          NAME     TYPE
    app_xxxxx   My App   app_store
    app_yyyyy   Web App  stripe
    ```

=== "CSV"

    ```bash
    rc apps list --output csv
    ```

    ```
    id,name,type
    app_xxxxx,My App,app_store
    app_yyyyy,Web App,stripe
    ```

---

## Common Workflows

### Pipe to jq

```bash
rc products list --output json | jq '.items[] | {id, identifier: .store_identifier}'
```

### Export to CSV

```bash
rc customers list --all --output csv > customers.csv
```

### Quick Lookup

```bash
rc entitlements get --entitlement-id entl_xxxxx --output minimal
```

---

## What's Next?

- Browse the [Commands](../commands/projects.md) reference for all 65+ commands
- Learn about [Output Formats](../reference/output-formats.md) in detail
- Set up [Shell Completion](../reference/shell-completion.md) for tab autocompletion

# Webhooks

Manage webhook endpoints for your RevenueCat project. Webhooks deliver real-time event notifications (purchases, renewals, cancellations, etc.) to your server.

## Available Commands

| Command | Description |
|---|---|
| `rc webhooks list` | List all configured webhooks |
| `rc webhooks create` | Create a new webhook endpoint |

---

## `rc webhooks list`

```bash
rc webhooks list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

### Example

```bash
rc webhooks list --output table
```

```
ID            URL
wh_xxxxx      https://api.example.com/webhooks/revenuecat
wh_yyyyy      https://staging.example.com/webhooks/rc
```

---

## `rc webhooks create`

Register a new webhook endpoint.

```bash
rc webhooks create --url https://api.example.com/webhooks/revenuecat
```

### Required Flags

| Flag | Description |
|---|---|
| `--url` | The HTTPS endpoint URL to receive webhook events |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

!!! note
    Webhook URLs must use HTTPS. RevenueCat will send a test event to verify the endpoint is reachable before activating it.

### Example

```bash
rc webhooks create \
  --url https://api.example.com/webhooks/revenuecat \
  --output pretty
```

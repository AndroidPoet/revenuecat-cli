# Paywalls

Manage paywalls in your RevenueCat project. Paywalls define the UI presentation layer for your offerings and are configured through the RevenueCat dashboard or API.

## Available Commands

| Command | Description |
|---|---|
| `rc paywalls list` | List all paywalls |
| `rc paywalls get` | Get a specific paywall |
| `rc paywalls create` | Create a new paywall |
| `rc paywalls delete` | Delete a paywall |

---

## `rc paywalls list`

```bash
rc paywalls list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

---

## `rc paywalls get`

```bash
rc paywalls get --paywall-id pw_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--paywall-id` | The paywall ID to retrieve |

---

## `rc paywalls create`

Create a new paywall and associate it with an offering.

```bash
rc paywalls create --offering-id ofr_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering to attach the paywall to |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

---

## `rc paywalls delete`

```bash
rc paywalls delete --paywall-id pw_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--paywall-id` | The paywall ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting a paywall removes it from its associated offering. Customers currently viewing this paywall will fall back to the default presentation.

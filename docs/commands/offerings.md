# Offerings

Manage offerings in your RevenueCat project. Offerings are the selection of products that are presented to a customer on your paywall. Each project can have multiple offerings, with one marked as "current."

## Available Commands

| Command | Description |
|---|---|
| `rc offerings list` | List all offerings |
| `rc offerings get` | Get a specific offering |
| `rc offerings create` | Create a new offering |
| `rc offerings update` | Update an offering |
| `rc offerings delete` | Delete an offering |

---

## `rc offerings list`

```bash
rc offerings list
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
rc offerings list --output table
```

```
ID            LOOKUP KEY     DISPLAY NAME      IS CURRENT
ofr_xxxxx     default        Default Offering  true
ofr_yyyyy     experiment_a   Experiment A      false
```

---

## `rc offerings get`

```bash
rc offerings get --offering-id ofr_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID to retrieve |

---

## `rc offerings create`

```bash
rc offerings create \
  --lookup-key default \
  --display-name "Default Offering"
```

### Required Flags

| Flag | Description |
|---|---|
| `--lookup-key` | Unique lookup key for the offering |
| `--display-name` | Human-readable display name |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--is-current` | Set as the current offering | `false` |
| `--metadata` | JSON metadata string | -- |
| `--output` | Output format | `json` |

### Example

```bash
rc offerings create \
  --lookup-key holiday_sale \
  --display-name "Holiday Sale" \
  --metadata '{"discount": "30%"}' \
  --is-current false
```

---

## `rc offerings update`

```bash
rc offerings update --offering-id ofr_xxxxx --display-name "Updated Name"
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID to update |

### Optional Flags

| Flag | Description |
|---|---|
| `--display-name` | New display name |
| `--is-current` | Set or unset as current offering |
| `--metadata` | Updated JSON metadata |

!!! tip
    Setting `--is-current true` on an offering automatically unsets the previously current offering.

---

## `rc offerings delete`

```bash
rc offerings delete --offering-id ofr_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting an offering also removes all its packages. Make sure no active paywalls reference this offering.

# Purchases

View and manage one-time purchases. Purchases represent non-recurring transactions such as lifetime unlocks or consumable items.

## Available Commands

| Command | Description |
|---|---|
| `rc purchases get` | Get details for a specific purchase |
| `rc purchases list-entitlements` | List entitlements granted by a purchase |
| `rc purchases refund` | Refund a purchase |

---

## `rc purchases get`

Retrieve full details for a purchase.

```bash
rc purchases get --purchase-id pur_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--purchase-id` | The purchase ID to retrieve |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

---

## `rc purchases list-entitlements`

List all entitlements granted by a specific purchase.

```bash
rc purchases list-entitlements --purchase-id pur_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--purchase-id` | The purchase ID |

---

## `rc purchases refund`

Issue a refund for a one-time purchase. This revokes any entitlements granted by the purchase.

```bash
rc purchases refund --purchase-id pur_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--purchase-id` | The purchase ID to refund |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Refunding a purchase immediately revokes entitlements and issues a refund through the original payment processor. This cannot be undone.

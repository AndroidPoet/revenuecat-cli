# Entitlements

Manage entitlements in your RevenueCat project. Entitlements represent levels of access that a customer can "unlock" through purchasing products.

## Available Commands

| Command | Description |
|---|---|
| `rc entitlements list` | List all entitlements |
| `rc entitlements get` | Get a specific entitlement |
| `rc entitlements create` | Create a new entitlement |
| `rc entitlements update` | Update an entitlement |
| `rc entitlements delete` | Delete an entitlement |
| `rc entitlements list-products` | List products attached to an entitlement |
| `rc entitlements attach-products` | Attach products to an entitlement |
| `rc entitlements detach-products` | Detach products from an entitlement |

---

## `rc entitlements list`

```bash
rc entitlements list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

---

## `rc entitlements get`

```bash
rc entitlements get --entitlement-id entl_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID to retrieve |

---

## `rc entitlements create`

```bash
rc entitlements create --lookup-key premium --display-name "Premium Access"
```

### Required Flags

| Flag | Description |
|---|---|
| `--lookup-key` | Unique lookup key for the entitlement |
| `--display-name` | Human-readable display name |

---

## `rc entitlements update`

```bash
rc entitlements update --entitlement-id entl_xxxxx --display-name "Pro Access"
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID to update |

### Optional Flags

| Flag | Description |
|---|---|
| `--display-name` | New display name |

---

## `rc entitlements delete`

```bash
rc entitlements delete --entitlement-id entl_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting an entitlement removes it from all offerings. Existing customer grants are not revoked.

---

## `rc entitlements list-products`

List all products attached to an entitlement.

```bash
rc entitlements list-products --entitlement-id entl_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID |

---

## `rc entitlements attach-products`

Attach one or more products to an entitlement. When a customer purchases any attached product, they receive this entitlement.

```bash
rc entitlements attach-products \
  --entitlement-id entl_xxxxx \
  --product-ids prod_aaaaa,prod_bbbbb
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID |
| `--product-ids` | Comma-separated list of product IDs to attach |

!!! tip
    You can attach products from different apps to the same entitlement. This is how cross-platform access works in RevenueCat.

---

## `rc entitlements detach-products`

Detach one or more products from an entitlement.

```bash
rc entitlements detach-products \
  --entitlement-id entl_xxxxx \
  --product-ids prod_aaaaa
```

### Required Flags

| Flag | Description |
|---|---|
| `--entitlement-id` | The entitlement ID |
| `--product-ids` | Comma-separated list of product IDs to detach |

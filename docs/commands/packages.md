# Packages

Manage packages within offerings. A package is a group of equivalent products across platforms within a single offering. For example, a "monthly" package might contain the iOS and Android versions of the same subscription.

## Available Commands

| Command | Description |
|---|---|
| `rc packages list` | List all packages in an offering |
| `rc packages get` | Get a specific package |
| `rc packages create` | Create a new package |
| `rc packages update` | Update a package |
| `rc packages delete` | Delete a package |
| `rc packages list-products` | List products attached to a package |
| `rc packages attach-products` | Attach products to a package |
| `rc packages detach-products` | Detach products from a package |

---

## `rc packages list`

List all packages in a specific offering.

```bash
rc packages list --offering-id ofr_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID to list packages for |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

---

## `rc packages get`

```bash
rc packages get --offering-id ofr_xxxxx --package-id pkg_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID to retrieve |

---

## `rc packages create`

```bash
rc packages create \
  --offering-id ofr_xxxxx \
  --lookup-key monthly \
  --display-name "Monthly Plan"
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering to create the package in |
| `--lookup-key` | Unique lookup key within the offering |
| `--display-name` | Human-readable display name |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

---

## `rc packages update`

```bash
rc packages update \
  --offering-id ofr_xxxxx \
  --package-id pkg_xxxxx \
  --display-name "Premium Monthly"
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID to update |

### Optional Flags

| Flag | Description |
|---|---|
| `--display-name` | New display name |

---

## `rc packages delete`

```bash
rc packages delete \
  --offering-id ofr_xxxxx \
  --package-id pkg_xxxxx \
  --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID to delete |
| `--confirm` | Skip confirmation prompt |

---

## `rc packages list-products`

List all products attached to a package.

```bash
rc packages list-products \
  --offering-id ofr_xxxxx \
  --package-id pkg_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID |

---

## `rc packages attach-products`

Attach products to a package. This defines which store products are included when a customer sees this package.

```bash
rc packages attach-products \
  --offering-id ofr_xxxxx \
  --package-id pkg_xxxxx \
  --product-ids prod_aaaaa,prod_bbbbb
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID |
| `--product-ids` | Comma-separated list of product IDs to attach |

!!! tip
    Attach one product per platform to a package. For example, attach both the iOS `app_store` product and the Android `play_store` product to the same "monthly" package for cross-platform parity.

---

## `rc packages detach-products`

Detach products from a package.

```bash
rc packages detach-products \
  --offering-id ofr_xxxxx \
  --package-id pkg_xxxxx \
  --product-ids prod_aaaaa
```

### Required Flags

| Flag | Description |
|---|---|
| `--offering-id` | The offering ID |
| `--package-id` | The package ID |
| `--product-ids` | Comma-separated list of product IDs to detach |

# Products

Manage products within your RevenueCat project. Products map to store-specific identifiers and define what your customers can purchase.

## Available Commands

| Command | Description |
|---|---|
| `rc products list` | List all products |
| `rc products get` | Get a specific product |
| `rc products create` | Create a new product |
| `rc products delete` | Delete a product |

---

## `rc products list`

List all products in the current project.

```bash
rc products list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--app-id` | Filter products by app | -- |
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

### Example

```bash
rc products list --app-id app_xxxxx --output table
```

```
ID             STORE IDENTIFIER        TYPE
prod_xxxxx     com.app.monthly         subscription
prod_yyyyy     com.app.lifetime        one_time
```

---

## `rc products get`

Retrieve details for a specific product.

```bash
rc products get --product-id prod_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--product-id` | The product ID to retrieve |

---

## `rc products create`

Create a new product.

```bash
rc products create \
  --app-id app_xxxxx \
  --store-identifier com.example.premium.monthly \
  --type subscription
```

### Required Flags

| Flag | Description |
|---|---|
| `--app-id` | The app this product belongs to |
| `--store-identifier` | The product identifier in the store |
| `--type` | Product type: `subscription` or `one_time` |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

### Examples

=== "Subscription"

    ```bash
    rc products create \
      --app-id app_xxxxx \
      --store-identifier com.example.premium.monthly \
      --type subscription
    ```

=== "One-Time Purchase"

    ```bash
    rc products create \
      --app-id app_xxxxx \
      --store-identifier com.example.lifetime \
      --type one_time
    ```

---

## `rc products delete`

Delete a product from the project.

```bash
rc products delete --product-id prod_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--product-id` | The product ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting a product detaches it from all entitlements and packages. Active subscriptions are not affected.

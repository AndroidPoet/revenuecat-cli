# Customers

Manage customers in your RevenueCat project. Customers represent the end users of your apps and hold subscriptions, purchases, entitlements, and attributes.

## Available Commands

| Command | Description |
|---|---|
| `rc customers list` | List all customers |
| `rc customers get` | Get a specific customer |
| `rc customers create` | Create a customer |
| `rc customers delete` | Delete a customer |
| `rc customers list-active-entitlements` | List a customer's active entitlements |
| `rc customers list-aliases` | List a customer's aliases |
| `rc customers list-attributes` | List a customer's attributes |
| `rc customers set-attributes` | Set attributes on a customer |
| `rc customers list-subscriptions` | List a customer's subscriptions |
| `rc customers list-purchases` | List a customer's purchases |
| `rc customers list-invoices` | List a customer's invoices |
| `rc customers transfer` | Transfer a customer to another ID |
| `rc customers grant-entitlement` | Grant a promotional entitlement |
| `rc customers revoke-entitlement` | Revoke a promotional entitlement |
| `rc customers assign-offering` | Assign a specific offering to a customer |

---

## `rc customers list`

```bash
rc customers list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

---

## `rc customers get`

```bash
rc customers get --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID to retrieve |

---

## `rc customers create`

Create a new customer with a specified ID.

```bash
rc customers create --customer-id my_custom_id
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The ID for the new customer |

---

## `rc customers delete`

```bash
rc customers delete --customer-id cust_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting a customer is permanent and removes all associated data including subscription history, attributes, and entitlements.

---

## `rc customers list-active-entitlements`

List entitlements currently active for a customer.

```bash
rc customers list-active-entitlements --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers list-aliases`

List all aliases associated with a customer.

```bash
rc customers list-aliases --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers list-attributes`

List all custom attributes set on a customer.

```bash
rc customers list-attributes --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers set-attributes`

Set custom attributes on a customer. Attributes are key-value pairs useful for segmentation and analytics.

```bash
rc customers set-attributes \
  --customer-id cust_xxxxx \
  --attributes '{"tier": "vip", "region": "us-west"}'
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |
| `--attributes` | JSON object of key-value pairs |

!!! note
    Setting an attribute key to an empty string (`""`) deletes that attribute.

---

## `rc customers list-subscriptions`

```bash
rc customers list-subscriptions --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers list-purchases`

```bash
rc customers list-purchases --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers list-invoices`

```bash
rc customers list-invoices --customer-id cust_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |

---

## `rc customers transfer`

Transfer a customer's data to a different customer ID.

```bash
rc customers transfer \
  --customer-id cust_xxxxx \
  --target-id cust_yyyyy
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The source customer ID |
| `--target-id` | The destination customer ID |

!!! warning
    Transfer merges subscription and purchase history into the target customer. The source customer will be left empty.

---

## `rc customers grant-entitlement`

Grant a promotional entitlement to a customer. This gives access without requiring a purchase.

```bash
rc customers grant-entitlement \
  --customer-id cust_xxxxx \
  --entitlement-id entl_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |
| `--entitlement-id` | The entitlement ID to grant |

---

## `rc customers revoke-entitlement`

Revoke a previously granted promotional entitlement.

```bash
rc customers revoke-entitlement \
  --customer-id cust_xxxxx \
  --entitlement-id entl_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |
| `--entitlement-id` | The entitlement ID to revoke |

---

## `rc customers assign-offering`

Override the offering shown to a specific customer. This is useful for A/B testing or granting special pricing.

```bash
rc customers assign-offering \
  --customer-id cust_xxxxx \
  --offering-id ofr_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--customer-id` | The customer ID |
| `--offering-id` | The offering ID to assign |

!!! tip
    To revert a customer to the default offering, assign the offering currently marked as `is_current: true`.

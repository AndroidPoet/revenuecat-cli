# Subscriptions

View and manage subscriptions. Subscriptions represent recurring purchases tied to a customer and product.

## Available Commands

| Command | Description |
|---|---|
| `rc subscriptions get` | Get details for a specific subscription |
| `rc subscriptions list-entitlements` | List entitlements granted by a subscription |
| `rc subscriptions cancel` | Cancel a subscription |
| `rc subscriptions refund` | Refund a subscription |

---

## `rc subscriptions get`

Retrieve full details for a subscription including its status, renewal date, and associated product.

```bash
rc subscriptions get --subscription-id sub_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--subscription-id` | The subscription ID to retrieve |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

---

## `rc subscriptions list-entitlements`

List all entitlements granted by a specific subscription.

```bash
rc subscriptions list-entitlements --subscription-id sub_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--subscription-id` | The subscription ID |

---

## `rc subscriptions cancel`

Cancel an active subscription. The customer retains access until the end of the current billing period.

```bash
rc subscriptions cancel --subscription-id sub_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--subscription-id` | The subscription ID to cancel |
| `--confirm` | Skip confirmation prompt |

!!! note
    Cancellation takes effect at the end of the current billing period. The customer is not immediately revoked.

---

## `rc subscriptions refund`

Issue a refund for a subscription. This revokes access immediately.

```bash
rc subscriptions refund --subscription-id sub_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--subscription-id` | The subscription ID to refund |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Refunding a subscription immediately revokes the customer's access and issues a refund through the original payment processor.

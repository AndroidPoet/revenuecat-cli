# Metrics

View key subscription metrics for your RevenueCat project. Metrics provide a real-time snapshot of your subscription business health.

## Available Commands

| Command | Description |
|---|---|
| `rc metrics overview` | Display key subscription metrics |

---

## `rc metrics overview`

Retrieve a high-level overview of your project's subscription metrics.

```bash
rc metrics overview
```

### Metrics Returned

| Metric | Description |
|---|---|
| **MRR** | Monthly Recurring Revenue |
| **Active Subscribers** | Number of customers with active paid subscriptions |
| **Active Trials** | Number of customers currently in a free trial |
| **Revenue** | Total revenue for the current period |

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |

### Examples

=== "Table Output"

    ```bash
    rc metrics overview --output table
    ```

    ```
    METRIC               VALUE
    MRR                  $12,450.00
    Active Subscribers   1,823
    Active Trials        342
    Revenue              $14,200.00
    ```

=== "JSON Output"

    ```bash
    rc metrics overview --output json
    ```

    ```json
    {
      "mrr": 12450.00,
      "active_subscribers": 1823,
      "active_trials": 342,
      "revenue": 14200.00
    }
    ```

=== "Minimal Output"

    ```bash
    rc metrics overview --output minimal
    ```

    ```
    12450.00 1823 342 14200.00
    ```

!!! tip
    Combine with `watch` for a live dashboard: `watch -n 60 rc metrics overview --output table`

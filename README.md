<div align="center">

<br>

<img src="assets/logo.png" alt="RevenueCat CLI" width="420">

<br>
<br>

**Your AI-powered subscription management companion for the command line.**

<br>

[![Release](https://img.shields.io/github/v/release/AndroidPoet/revenuecat-cli?style=for-the-badge&color=4B48F2&label=Latest)](https://github.com/AndroidPoet/revenuecat-cli/releases/latest)
&nbsp;
[![Downloads](https://img.shields.io/github/downloads/AndroidPoet/revenuecat-cli/total?style=for-the-badge&color=E8514A)](https://github.com/AndroidPoet/revenuecat-cli/releases)
&nbsp;
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
&nbsp;
[![License](https://img.shields.io/badge/License-MIT-F5A623?style=for-the-badge)](LICENSE)

</div>

<br>

## Install

```bash
brew tap AndroidPoet/tap && brew install revenuecat-cli
```

Or download from [Releases](https://github.com/AndroidPoet/revenuecat-cli/releases/latest). After install, use `revenuecat-cli` or the alias `rc`.

## Command Map

```
rc
├── projects
│   ├── list                          # List all projects
│   └── create                        # Create a new project
├── apps
│   ├── list / get / create           # CRUD operations
│   ├── update / delete               # Modify or remove
│   └── api-keys                      # List public API keys
├── products
│   ├── list / get / create           # Manage store products
│   └── delete                        # Remove a product
├── entitlements
│   ├── list / get / create           # Manage entitlements
│   ├── update / delete               # Modify or remove
│   └── list-products / attach / detach  # Product associations
├── offerings
│   ├── list / get / create           # Manage offerings
│   └── update / delete               # Modify or remove
├── packages
│   ├── list / get / create           # Manage packages
│   ├── update / delete               # Modify or remove
│   └── list-products / attach / detach  # Product associations
├── customers
│   ├── list / get / create / delete  # CRUD operations
│   ├── list-subscriptions            # View subscriptions
│   ├── list-purchases                # View purchases
│   ├── list-invoices                 # View invoices
│   ├── list-active-entitlements      # View active entitlements
│   ├── list-aliases / list-attributes   # View customer data
│   ├── set-attributes               # Update attributes
│   ├── grant-entitlement             # Grant access
│   ├── revoke-entitlement            # Revoke access
│   ├── assign-offering               # Override offering
│   └── transfer                      # Transfer to another customer
├── subscriptions
│   ├── get / list-entitlements       # View details
│   ├── cancel                        # Cancel subscription
│   └── refund                        # Refund subscription
├── purchases
│   ├── get / list-entitlements       # View details
│   └── refund                        # Refund purchase
├── paywalls
│   ├── list / get / create           # Manage paywalls
│   └── delete                        # Remove a paywall
├── metrics
│   └── overview                      # MRR, subscribers, trials, revenue
├── webhooks
│   ├── list                          # View webhook integrations
│   └── create                        # Create a webhook
├── audit-logs
│   └── list                          # View audit log entries
└── auth
    ├── login / switch                # Manage profiles
    ├── list / current                # View profiles
    └── delete                        # Remove a profile
```

## Setup

```bash
rc auth login --api-key sk_your_key_here
rc init --project proj_your_project_id
rc doctor
```

## Commands

### Projects

```bash
rc projects list
rc projects create --name "My Project"
```

### Apps

```bash
rc apps list
rc apps get --app-id app_xxx
rc apps create --name "My App" --type play_store
rc apps update --app-id app_xxx --name "New Name"
rc apps delete --app-id app_xxx --confirm
rc apps api-keys --app-id app_xxx
```

### Products

```bash
rc products list
rc products get --product-id prod_xxx
rc products create --store-identifier com.app.monthly --type subscription --app-id app_xxx
rc products delete --product-id prod_xxx --confirm
```

### Entitlements

```bash
rc entitlements list
rc entitlements get --entitlement-id entl_xxx
rc entitlements create --lookup-key premium --display-name "Premium"
rc entitlements update --entitlement-id entl_xxx --display-name "Premium+"
rc entitlements delete --entitlement-id entl_xxx --confirm
rc entitlements list-products --entitlement-id entl_xxx
rc entitlements attach-products --entitlement-id entl_xxx --product-ids prod1,prod2
rc entitlements detach-products --entitlement-id entl_xxx --product-ids prod1
```

### Offerings

```bash
rc offerings list
rc offerings get --offering-id ofrngs_xxx
rc offerings create --lookup-key default --display-name "Default Offering"
rc offerings update --offering-id ofrngs_xxx --is-current
rc offerings update --offering-id ofrngs_xxx --metadata '{"theme":"dark"}'
rc offerings delete --offering-id ofrngs_xxx --confirm
```

### Packages

```bash
rc packages list --offering-id ofrngs_xxx
rc packages get --package-id pkg_xxx
rc packages create --offering-id ofrngs_xxx --lookup-key monthly --display-name "Monthly"
rc packages update --package-id pkg_xxx --display-name "Monthly Plan"
rc packages delete --package-id pkg_xxx --confirm
rc packages list-products --package-id pkg_xxx
rc packages attach-products --package-id pkg_xxx --product-ids prod1,prod2
rc packages detach-products --package-id pkg_xxx --product-ids prod1
```

### Customers

```bash
rc customers list
rc customers get --customer-id user_xxx
rc customers create --customer-id user_xxx
rc customers delete --customer-id user_xxx --confirm
rc customers list-active-entitlements --customer-id user_xxx
rc customers list-subscriptions --customer-id user_xxx
rc customers list-purchases --customer-id user_xxx
rc customers list-invoices --customer-id user_xxx
rc customers list-aliases --customer-id user_xxx
rc customers list-attributes --customer-id user_xxx
rc customers set-attributes --customer-id user_xxx --attributes '{"key":"value"}'
rc customers transfer --customer-id user_xxx --target-id user_yyy
rc customers grant-entitlement --customer-id user_xxx --entitlement-id entl_xxx
rc customers revoke-entitlement --customer-id user_xxx --entitlement-id entl_xxx
rc customers assign-offering --customer-id user_xxx --offering-id ofrngs_xxx
```

### Subscriptions

```bash
rc subscriptions get --subscription-id sub_xxx
rc subscriptions list-entitlements --subscription-id sub_xxx
rc subscriptions cancel --subscription-id sub_xxx --confirm
rc subscriptions refund --subscription-id sub_xxx --confirm
```

### Purchases

```bash
rc purchases get --purchase-id pur_xxx
rc purchases list-entitlements --purchase-id pur_xxx
rc purchases refund --purchase-id pur_xxx --confirm
```

### Paywalls

```bash
rc paywalls list
rc paywalls get --paywall-id pw_xxx
rc paywalls create --offering-id ofrngs_xxx
rc paywalls delete --paywall-id pw_xxx --confirm
```

### Metrics

```bash
rc metrics overview
```

### Webhooks

```bash
rc webhooks list
rc webhooks create --url https://example.com/webhook
```

### Audit Logs

```bash
rc audit-logs list
```

### Auth

```bash
rc auth login --api-key sk_xxx --name production
rc auth login --api-key sk_xxx --name staging --default-project proj_staging
rc auth switch --name staging
rc auth list
rc auth current
rc auth delete --name old --confirm
```

### Utilities

```bash
rc doctor
rc init --project proj_xxx
rc completion zsh > "${fpath[1]}/_rc"
rc version
```

## Output Formats

```bash
rc apps list                # JSON (default)
rc apps list --pretty       # Pretty JSON
rc apps list -o table       # Table
rc apps list -o csv         # CSV
rc apps list -o tsv         # TSV
rc apps list -o yaml        # YAML
rc apps list -o minimal     # IDs only
```

## Pagination

```bash
rc products list --limit 50
rc products list --starting-after prod_xxx
rc products list --all
```

## Contributing

```bash
make build
make test
make lint
```

## License

MIT

---

<div align="center">
<sub>Not affiliated with RevenueCat Inc.</sub>
</div>

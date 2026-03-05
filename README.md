<div align="center">

<br>

<img src="assets/logo.png" alt="RevenueCat CLI" width="420">

<br>
<br>

**Your AI-powered subscription management companion for the command line.**

<br>

[![Release](https://img.shields.io/github/v/release/AndroidPoet/revenuecat-cli?style=for-the-badge&color=4B48F2&label=Latest)](https://github.com/AndroidPoet/revenuecat-cli/releases/latest)
&nbsp;
[![Downloads](https://img.shields.io/github/downloads/AndroidPoet/revenuecat-cli/total?style=for-the-badge&color=E8514A&v=1)](https://github.com/AndroidPoet/revenuecat-cli/releases)
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

## Command Overview

> **70+ commands** across **14 resource groups** covering the full RevenueCat API v2.

| Category | Commands | What you can do |
|:---------|:---------|:----------------|
| **Projects** | `list` `create` | Manage your RevenueCat projects |
| **Apps** | `list` `get` `create` `update` `delete` `api-keys` | Configure app store connections |
| **Products** | `list` `get` `create` `delete` | Define subscription and one-time products |
| **Entitlements** | `list` `get` `create` `update` `delete` `attach-products` `detach-products` | Control access to premium features |
| **Offerings** | `list` `get` `create` `update` `delete` | Group packages for remote configuration |
| **Packages** | `list` `get` `create` `update` `delete` `attach-products` `detach-products` | Bundle products within offerings |
| **Customers** | `list` `get` `create` `delete` + **11 more** | Full customer lifecycle management |
| **Subscriptions** | `get` `list-entitlements` `cancel` `refund` | Manage active subscriptions |
| **Purchases** | `get` `list-entitlements` `refund` | View and refund purchases |
| **Paywalls** | `list` `get` `create` `delete` | Manage paywall configurations |
| **Metrics** | `overview` | MRR, active subscribers, trials, revenue |
| **Webhooks** | `list` `create` | Set up webhook integrations |
| **Audit Logs** | `list` | Track changes and access history |
| **Auth** | `login` `switch` `list` `current` `delete` | Manage API key profiles |

## New in v0.3.0

| Feature | Command | Description |
|:--------|:--------|:------------|
| **Status Dashboard** | `rc status` | One-command overview of your entire project |
| **Live Watch** | `rc watch metrics` | Auto-refreshing terminal metrics dashboard |
| **Project Diff** | `rc diff --source A --target B` | Compare entitlements and offerings between projects |
| **Export/Import** | `rc export` / `rc import` | Backup and migrate project configuration as YAML |
| **Full Report** | `rc report` | Export entire project data as HTML (PDF-ready), JSON, or YAML |
| **Dynamic Completion** | Tab on `--app-id`, `--product-id`, etc. | Live API-powered shell completions |
| **Colored Output** | Automatic | Green checkmarks, cyan info, red errors |
| **CI Pipeline** | GitHub Actions | Build, test, lint on every push |

## Setup

> Get your API key at [app.revenuecat.com/settings/api-keys](https://app.revenuecat.com/settings/api-keys)

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

### Status & Watch

```bash
rc status                                       # Project dashboard
rc watch metrics                                # Live metrics (Ctrl+C to stop)
rc watch metrics --interval 10s                 # Custom refresh interval
```

### Export & Import

```bash
rc export --file my-project.yaml                # Export config to YAML
rc import --file my-project.yaml --confirm      # Import config from YAML
rc import --file my-project.yaml --dry-run      # Preview import
```

### Report

```bash
rc report                                       # HTML report (open in browser, print to PDF)
rc report --format json --file report.json      # Full data as JSON
rc report --format yaml --file report.yaml      # Full data as YAML
```

### Diff

```bash
rc diff --source proj_staging --target proj_prod  # Compare two projects
```

### Utilities

```bash
rc doctor
rc init --project proj_xxx
rc completion zsh > "${fpath[1]}/_rc"
rc version
```

## Agent Skills

Use `rc` with AI coding agents. Install the skill pack and your agent learns every command — products, entitlements, offerings, customers, metrics, and more.

```bash
npx skills add AndroidPoet/revenuecat-cli-skills
```

Then just ask:

```
Create a premium entitlement and attach my monthly subscription product
```
```
Show me the current MRR and active subscriber count
```
```
Look up customer user_123 and list their active entitlements
```

[Browse all 8 skills →](https://github.com/AndroidPoet/revenuecat-cli-skills)

---

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

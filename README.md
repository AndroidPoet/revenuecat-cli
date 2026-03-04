<div align="center">

<br>

<img src="assets/logo.png" alt="RevenueCat CLI" width="420">

<br>
<br>

**Manage in-app subscriptions from your terminal.**

<br>

[![Release](https://img.shields.io/github/v/release/AndroidPoet/revenuecat-cli?style=for-the-badge&color=4B48F2&label=Latest)](https://github.com/AndroidPoet/revenuecat-cli/releases/latest)
&nbsp;
[![Downloads](https://img.shields.io/github/downloads/AndroidPoet/revenuecat-cli/total?style=for-the-badge&color=E8514A)](https://github.com/AndroidPoet/revenuecat-cli/releases)
&nbsp;
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
&nbsp;
[![MIT](https://img.shields.io/badge/License-MIT-F5A623?style=for-the-badge)](LICENSE)

</div>

<br>

## Install

```bash
brew tap AndroidPoet/tap && brew install revenuecat-cli
```

Or download from [Releases](https://github.com/AndroidPoet/revenuecat-cli/releases/latest). After install, use `revenuecat-cli` or the alias `rc`.

## Setup

```bash
# 1. Login with your RevenueCat API v2 secret key
rc auth login --api-key sk_your_key_here

# 2. Set your project
rc init --project proj_your_project_id

# 3. Verify
rc doctor
```

## Commands

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
rc products create --store-identifier com.app.monthly --type subscription --app-id app_xxx
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
rc offerings create --lookup-key default --display-name "Default Offering"
rc offerings update --offering-id ofrngs_xxx --is-current
rc offerings update --offering-id ofrngs_xxx --metadata '{"theme":"dark"}'
```

### Packages

```bash
rc packages list --offering-id ofrngs_xxx
rc packages create --offering-id ofrngs_xxx --lookup-key monthly --display-name "Monthly"
rc packages attach-products --package-id pkg_xxx --product-ids prod1,prod2
rc packages detach-products --package-id pkg_xxx --product-ids prod1
```

### Customers

```bash
rc customers get --customer-id user_xxx
rc customers delete --customer-id user_xxx --confirm
```

### Paywalls

```bash
rc paywalls create --offering-id ofrngs_xxx
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

## Environment Variables

| Variable | Description |
|:---------|:------------|
| `RC_API_KEY` | API v2 secret key |
| `RC_PROJECT` | Project ID |
| `RC_PROFILE` | Auth profile |
| `RC_OUTPUT` | Output format |
| `RC_DEBUG` | Log HTTP requests |
| `RC_TIMEOUT` | Request timeout |

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

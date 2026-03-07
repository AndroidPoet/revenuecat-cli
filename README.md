<div align="center">

<br>

<img src="assets/logo.png" alt="RevenueCat CLI" width="420">

<br>
<br>

**The command-line tool for RevenueCat — manage subscriptions, analyze charts, and export reports from your terminal.**

<br>

[![Release](https://img.shields.io/github/v/release/AndroidPoet/revenuecat-cli?style=for-the-badge&color=4B48F2&label=Latest)](https://github.com/AndroidPoet/revenuecat-cli/releases/latest)
&nbsp;
[![Downloads](https://img.shields.io/github/downloads/AndroidPoet/revenuecat-cli/total?style=for-the-badge&color=E8514A&v=2)](https://github.com/AndroidPoet/revenuecat-cli/releases)
&nbsp;
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
&nbsp;
[![License](https://img.shields.io/badge/License-MIT-F5A623?style=for-the-badge)](LICENSE)

</div>

<br>

## Quick Install

```bash
brew tap AndroidPoet/tap && brew install revenuecat-cli
```

Or download from [Releases](https://github.com/AndroidPoet/revenuecat-cli/releases/latest). After install, use `revenuecat-cli` or the alias `rc`. See [all install options](#install-1) below.

## What's New

### v0.5.0 — Magic Setup

| Feature | Command | Description |
|:--------|:--------|:------------|
| **Magic Setup** | `rc magicsetup` | One-click interactive setup for iOS, Android, or both |
| **5 Templates** | Freemium, Paywall, Trial, Tiered, Consumable | Preset configurations ready to customize |
| **Price Input** | Per-product pricing | Set display prices during setup |
| **Dry Run** | `rc magicsetup --dry-run` | Preview everything before creating |

### v0.4.0 — Charts & Analytics

| Feature | Command | Description |
|:--------|:--------|:------------|
| **21 Chart Types** | `rc charts list` | Revenue, MRR, ARR, churn, trials, LTV, retention, and more |
| **Chart Data** | `rc charts get revenue` | Fetch time-series data with resolution, date range, currency, segments |
| **Chart Options** | `rc charts options revenue` | Discover available filters, segments, and resolutions per chart |
| **Visual Export** | `rc charts export --format pdf` | SVG chart reports as HTML, PDF, JSON, YAML, or CSV |

### v0.3.0 — Operations & Monitoring

| Feature | Command | Description |
|:--------|:--------|:------------|
| **Status Dashboard** | `rc status` | One-command overview of your entire project |
| **Live Watch** | `rc watch metrics` | Auto-refreshing terminal metrics dashboard |
| **Project Diff** | `rc diff --source A --target B` | Compare entitlements and offerings between projects |
| **Export/Import** | `rc export` / `rc import` | Backup and migrate project configuration as YAML |
| **Full Report** | `rc report --format pdf` | Export entire project data as HTML, PDF, JSON, or YAML |
| **Dynamic Completion** | Tab on `--app-id`, `--product-id`, etc. | Live API-powered shell completions |

## Setup

> Get your API key at [app.revenuecat.com/settings/api-keys](https://app.revenuecat.com/settings/api-keys)

```bash
rc auth login --api-key sk_your_key_here
rc init --project proj_your_project_id
rc doctor
```

## Command Overview

> **80+ commands** across **16 resource groups** covering the full RevenueCat API v2.

| Category | Commands | What you can do |
|:---------|:---------|:----------------|
| **Magic Setup** | `magicsetup` | One-click setup wizard for iOS, Android, or both |
|:---------|:---------|:----------------|
| **Charts** | `list` `get` `options` `export` | 21 chart types with SVG visual export |
| **Projects** | `list` `create` | Manage your RevenueCat projects |
| **Apps** | `list` `get` `create` `update` `delete` `api-keys` | Configure app store connections |
| **Products** | `list` `get` `create` `delete` | Define subscription and one-time products |
| **Entitlements** | `list` `get` `create` `update` `delete` `attach` `detach` | Control access to premium features |
| **Offerings** | `list` `get` `create` `update` `delete` | Group packages for remote configuration |
| **Packages** | `list` `get` `create` `update` `delete` `attach` `detach` | Bundle products within offerings |
| **Customers** | `list` `get` `create` `delete` + **11 more** | Full customer lifecycle management |
| **Subscriptions** | `get` `list-entitlements` `cancel` `refund` | Manage active subscriptions |
| **Purchases** | `get` `list-entitlements` `refund` | View and refund purchases |
| **Paywalls** | `list` `get` `create` `delete` | Manage paywall configurations |
| **Metrics** | `overview` | MRR, active subscribers, trials, revenue |
| **Webhooks** | `list` `create` | Set up webhook integrations |
| **Audit Logs** | `list` | Track changes and access history |
| **Auth** | `login` `switch` `list` `current` `delete` | Manage API key profiles |

## Magic Setup

Set up your entire RevenueCat offering stack in one command — apps, products, entitlements, offerings, and packages, all wired together.

```bash
rc magicsetup
```

### What it does

```
  ┌─────────────────────────────────────────────┐
  │            Magic Setup                      │
  │                                             │
  │  Apps  Products  Entitlements  Offerings    │
  │  All wired together, ready to go.           │
  └─────────────────────────────────────────────┘
```

**1. Pick your platform**
```
  1) iOS only           (App Store)
  2) Android only       (Play Store)
  3) Both iOS + Android  (recommended)
```

**2. Choose a template (or build custom)**
```
  1) Freemium          Weekly, Monthly, Annual + Lifetime
  2) Simple Paywall    Monthly + Annual
  3) Trial-First       Free trial then Monthly / Annual
  4) Tiered            Basic + Pro tiers with Monthly / Annual
  5) Consumable        Credit / coin packs
  6) Custom            Build your own from scratch
```

**3. Set prices and customize**
```
  Package                Type           Price          Product IDs
  ────────────────────── ────────────── ────────────── ─────────────────────────
  Weekly                 subscription   $1.99/week     iOS + Android
  Monthly                subscription   $9.99/month    iOS + Android
  Annual                 subscription   $49.99/year    iOS + Android
  Lifetime               one_time       $149.99        iOS + Android
```

**4. Confirm and it creates everything**
```
✓ iOS app: app_1a2b3c
✓ Android app: app_4d5e6f
✓ 8 products created
✓ Entitlements attached
✓ Offering set as current
✓ 4 packages wired up
```

### Options

```bash
rc magicsetup              # Interactive wizard
rc magicsetup --dry-run    # Preview without creating anything
```

## Charts & Analytics

Fetch, explore, and export visual analytics from all 21 RevenueCat chart types.

```bash
rc charts list                                  # List all 21 chart types
rc charts get revenue --resolution month        # Fetch revenue data (monthly)
rc charts get mrr --start-date 2024-01-01 --end-date 2024-12-31
rc charts get churn --currency EUR --segment country
rc charts options revenue                       # Discover filters, segments, resolutions
```

### Export visual reports

```bash
rc charts export --format html                  # SVG chart report (open in browser)
rc charts export --format pdf                   # PDF via Chrome headless
rc charts export --format csv                   # Flat CSV for spreadsheets
rc charts export --format json                  # Structured JSON
rc charts export --format yaml                  # YAML
rc charts export --charts revenue,mrr,actives   # Custom chart selection
```

**Available charts:** `revenue` `mrr` `arr` `mrr_movement` `actives` `actives_movement` `actives_new` `trials` `trials_movement` `trials_new` `trial_conversion_rate` `conversion_to_paying` `churn` `refund_rate` `subscription_retention` `subscription_status` `customers_active` `customers_new` `ltv_per_customer` `ltv_per_paying_customer` `cohort_explorer`

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
rc packages attach-products --package-id pkg_xxx --product-ids prod1,prod2
rc packages detach-products --package-id pkg_xxx --product-ids prod1
```

### Customers

```bash
rc customers get --customer-id user_xxx
rc customers create --customer-id user_xxx
rc customers delete --customer-id user_xxx --confirm
rc customers list-active-entitlements --customer-id user_xxx
rc customers list-subscriptions --customer-id user_xxx
rc customers list-purchases --customer-id user_xxx
rc customers list-invoices --customer-id user_xxx
rc customers set-attributes --customer-id user_xxx --attributes '{"key":"value"}'
rc customers transfer --customer-id user_xxx --target-id user_yyy
rc customers grant-entitlement --customer-id user_xxx --entitlement-id entl_xxx
rc customers revoke-entitlement --customer-id user_xxx --entitlement-id entl_xxx
```

### Subscriptions & Purchases

```bash
rc subscriptions get --subscription-id sub_xxx
rc subscriptions cancel --subscription-id sub_xxx --confirm
rc subscriptions refund --subscription-id sub_xxx --confirm
rc purchases get --purchase-id pur_xxx
rc purchases refund --purchase-id pur_xxx --confirm
```

### Paywalls

```bash
rc paywalls list
rc paywalls get --paywall-id pw_xxx
rc paywalls create --offering-id ofrngs_xxx
rc paywalls delete --paywall-id pw_xxx --confirm
```

### Metrics & Status

```bash
rc metrics overview                             # MRR, subscribers, trials, revenue
rc status                                       # Full project dashboard
rc watch metrics                                # Live auto-refreshing metrics
rc watch metrics --interval 10s                 # Custom refresh interval
```

### Reports

```bash
rc report                                       # HTML project report
rc report --format pdf                          # Direct PDF export
rc report --format json --file report.json      # Structured JSON
rc report --format yaml --file report.yaml      # YAML
```

### Export, Import & Diff

```bash
rc export --file my-project.yaml                # Export config to YAML
rc import --file my-project.yaml --dry-run      # Preview import
rc import --file my-project.yaml --confirm      # Import config from YAML
rc diff --source proj_staging --target proj_prod  # Compare two projects
```

### Auth & Setup

```bash
rc auth login --api-key sk_xxx --name production
rc auth login --api-key sk_xxx --name staging --default-project proj_staging
rc auth switch --name staging
rc auth list
rc auth current
rc auth delete --name old --confirm
rc doctor                                       # Verify configuration
rc init --project proj_xxx                      # Initialize project config
rc completion zsh > "${fpath[1]}/_rc"           # Shell completions
```

### Webhooks & Audit

```bash
rc webhooks list
rc webhooks create --url https://example.com/webhook
rc audit-logs list
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

## Agent Skills

Use `rc` with AI coding agents. Install the skill pack and your agent learns every command.

```bash
npx skills add AndroidPoet/revenuecat-cli-skills
```

Then just ask:

```
Export my revenue and MRR charts as a PDF report
```
```
Show me the current MRR and active subscriber count
```
```
Compare my staging and production project configurations
```

[Browse all 10 skills →](https://github.com/AndroidPoet/revenuecat-cli-skills)

## Install

### Homebrew (macOS / Linux)

```bash
brew tap AndroidPoet/tap && brew install revenuecat-cli
```

### Go Install

```bash
go install github.com/AndroidPoet/revenuecat-cli/cmd/revenuecat-cli@latest
```

### Binary Download

Download the latest binary for your platform from [Releases](https://github.com/AndroidPoet/revenuecat-cli/releases/latest).

| Platform | Architecture | File |
|:---------|:-------------|:-----|
| macOS | Apple Silicon | `revenuecat-cli_*_darwin_arm64.tar.gz` |
| macOS | Intel | `revenuecat-cli_*_darwin_amd64.tar.gz` |
| Linux | x86_64 | `revenuecat-cli_*_linux_amd64.tar.gz` |
| Linux | ARM64 | `revenuecat-cli_*_linux_arm64.tar.gz` |
| Windows | x86_64 | `revenuecat-cli_*_windows_amd64.zip` |

### Linux Packages

```bash
# Debian / Ubuntu
sudo dpkg -i revenuecat-cli_*_amd64.deb

# RHEL / Fedora
sudo rpm -i revenuecat-cli_*_amd64.rpm
```

## How Releases Work

Releases are fully automated via GitHub Actions and [GoReleaser](https://goreleaser.com).

```
Push v* tag  →  CI runs tests  →  GoReleaser builds  →  GitHub Release created
                                       │
                                       ├── Binaries (macOS, Linux, Windows × amd64/arm64)
                                       ├── .deb and .rpm packages
                                       ├── Homebrew formula updated (AndroidPoet/homebrew-tap)
                                       └── Checksums (checksums.txt)
```

### Creating a release

```bash
# Tag and push
git tag v0.5.0
git push origin v0.5.0
```

That's it. The [release workflow](.github/workflows/release.yml) handles:
1. Running tests
2. Building binaries for all platforms (CGO_ENABLED=0)
3. Creating `.tar.gz` / `.zip` archives
4. Generating `.deb` and `.rpm` packages
5. Publishing the GitHub Release
6. Updating the Homebrew formula in [AndroidPoet/homebrew-tap](https://github.com/AndroidPoet/homebrew-tap)

### CI

Every push and PR to `master` runs:
- `go build ./...`
- `go test ./... -v -race`
- `go vet ./...`
- `golangci-lint`

## Contributing

```bash
git clone https://github.com/AndroidPoet/revenuecat-cli.git
cd revenuecat-cli
make deps      # Download dependencies
make build     # Build binary to bin/
make test      # Run tests
make lint      # Run linter
```

## License

MIT

---

<div align="center">
<sub>Not affiliated with RevenueCat Inc.</sub>
</div>

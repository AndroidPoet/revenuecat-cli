---
hide:
  - navigation
  - toc
---

# RevenueCat CLI

<div class="hero" markdown>

**Your AI-powered subscription management companion for the command line.**

Manage your entire RevenueCat project without leaving the terminal. From products and entitlements to customers and metrics -- everything at your fingertips.

</div>

---

## Why RevenueCat CLI?

| Feature | Details |
|---|---|
| **65+ Commands** | Full coverage of the RevenueCat REST API v2 |
| **14 Resource Groups** | Projects, Apps, Products, Entitlements, Offerings, Packages, Customers, Subscriptions, Purchases, Paywalls, Metrics, Webhooks, Audit Logs, Auth |
| **7 Output Formats** | JSON, Pretty JSON, Table, CSV, TSV, YAML, Minimal |
| **Cursor-Based Pagination** | Efficiently traverse large datasets with `--limit`, `--starting-after`, and `--all` |
| **Multiple Auth Profiles** | Switch between accounts and projects seamlessly |
| **Shell Completion** | Bash, Zsh, Fish, and PowerShell support |

---

## Quick Install

```bash
brew tap AndroidPoet/tap && brew install revenuecat-cli
```

Verify the installation:

```bash
rc version
```

---

## Get Started in 60 Seconds

```bash
# 1. Authenticate
rc auth login --api-key sk_xxxxx --name production

# 2. Set your project
rc init --project proj_xxxxx

# 3. Verify everything works
rc doctor

# 4. Start managing
rc apps list
rc products list --output table
```

---

## What's Next?

<div class="grid cards" markdown>

- [**Installation**](getting-started/installation.md) -- Install via Homebrew or direct download
- [**Configuration**](getting-started/configuration.md) -- Set up auth profiles and project defaults
- [**Quick Start**](getting-started/quickstart.md) -- Your first commands in under a minute
- [**Commands**](commands/projects.md) -- Browse all 65+ commands across 14 resource groups

</div>

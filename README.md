<div align="center">

<br>

<img src="assets/logo.png" alt="RevenueCat CLI" width="420">

<br>
<br>

<h3>Subscriptions belong in the terminal, not the dashboard.</h3>

<br>

[![Release](https://img.shields.io/github/v/release/AndroidPoet/revenuecat-cli?style=for-the-badge&color=4B48F2&label=Latest)](https://github.com/AndroidPoet/revenuecat-cli/releases/latest)
&nbsp;
[![Downloads](https://img.shields.io/github/downloads/AndroidPoet/revenuecat-cli/total?style=for-the-badge&color=E8514A)](https://github.com/AndroidPoet/revenuecat-cli/releases)
&nbsp;
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
&nbsp;
[![MIT](https://img.shields.io/badge/License-MIT-F5A623?style=for-the-badge)](LICENSE)

<br>

```
$ rc offerings list -o table

LOOKUP_KEY    DISPLAY_NAME        IS_CURRENT
─             ─                   ─
default       Default Offering    true
premium       Premium Tier        false
enterprise    Enterprise Plan     false
```

<br>

[Install](#-install) · [Quick Start](#-quick-start) · [All Commands](#-all-commands) · [Scripting](#-scripting) · [Config](#-configuration)

<br>

</div>

---

<br>

## 🎯 Why

| 😤 Without `rc` | ⚡ With `rc` |
|:---|:---|
| Open dashboard → Settings → Apps → scroll... | `rc apps list` |
| Click through 4 screens to create a product | `rc products create --store-identifier com.app.pro --type subscription --app-id app_xxx` |
| Manually wire products to entitlements | `rc entitlements attach-products --product-ids prod1,prod2` |
| "Who changed the current offering?" | `rc offerings list -o table` |
| Export product catalog? Screenshot the UI... | `rc products list --all -o csv > products.csv` |
| Debug a customer's subscription state | `rc customers get --customer-id user_12345 --pretty` |

<br>

---

<br>

## 📦 Install

<table>
<tr>
<td>

**Homebrew**

```bash
brew tap AndroidPoet/tap
brew install revenuecat-cli
```

</td>
<td>

**Script**

```bash
curl -fsSL https://raw.githubusercontent.com/AndroidPoet/revenuecat-cli/main/install.sh | bash
```

</td>
<td>

**Source**

```bash
git clone https://github.com/AndroidPoet/revenuecat-cli.git
cd revenuecat-cli && make build
```

</td>
</tr>
</table>

> Use `revenuecat-cli` or the alias **`rc`**. Same binary, your choice.

<br>

---

<br>

## 🚀 Quick Start

**1.** Grab your **v2 secret key** from [RevenueCat → API Keys](https://app.revenuecat.com/settings/api-keys)

**2.** Login and set your project:

```bash
rc auth login --api-key sk_your_key_here
rc init --project proj_your_project_id
```

**3.** Verify everything:

```bash
$ rc doctor

✓ Configuration: OK (/Users/you/.revenuecat-cli/config.json)
✓ API Key: OK (sk_tes...here)
✓ Project ID: OK (proj_your_project_id)
✓ API Connectivity: OK (Successfully connected to RevenueCat API)

All checks passed!
```

**4.** Go:

```bash
rc apps list
rc entitlements list -o table
rc offerings list --pretty
```

<br>

---

<br>

## 📖 All Commands

<details open>
<summary><h3>🏗️ Apps — 6 commands</h3></summary>

> Create and manage apps across **iOS, Android, Stripe, Amazon, Mac, Roku, and Web**.

```bash
rc apps list                                            # List all apps
rc apps get --app-id app_xxx                            # Get details
rc apps create --name "My App" --type play_store        # Create
rc apps update --app-id app_xxx --name "New Name"       # Rename
rc apps delete --app-id app_xxx --confirm               # Delete
rc apps api-keys --app-id app_xxx                       # Public API keys
```

**App types:** `app_store` · `play_store` · `stripe` · `amazon` · `mac_app_store` · `roku` · `web`

</details>

<details open>
<summary><h3>📦 Products — 2 commands</h3></summary>

```bash
rc products list                                        # All products
rc products create --store-identifier com.app.monthly \
  --type subscription --app-id app_xxx                  # Create
```

**Product types:** `subscription` · `one_time`

</details>

<details open>
<summary><h3>🔑 Entitlements — 8 commands</h3></summary>

> The heart of RevenueCat. Full CRUD + product association management.

```bash
rc entitlements list                                    # List all
rc entitlements get --entitlement-id entl_xxx            # Get one
rc entitlements create --lookup-key premium \
  --display-name "Premium"                              # Create
rc entitlements update --entitlement-id entl_xxx \
  --display-name "Premium+"                             # Rename
rc entitlements delete --entitlement-id entl_xxx \
  --confirm                                             # Delete
```

**Product wiring:**

```bash
rc entitlements list-products --entitlement-id entl_xxx                        # See what's attached
rc entitlements attach-products --entitlement-id entl_xxx --product-ids a,b    # Attach
rc entitlements detach-products --entitlement-id entl_xxx --product-ids a      # Detach
```

</details>

<details open>
<summary><h3>🎁 Offerings — 3 commands</h3></summary>

```bash
rc offerings list                                       # All offerings
rc offerings create --lookup-key default \
  --display-name "Default Offering"                     # Create
rc offerings update --offering-id ofrngs_xxx \
  --is-current                                          # Make current
rc offerings update --offering-id ofrngs_xxx \
  --metadata '{"theme":"dark","version":2}'             # Set metadata
```

</details>

<details open>
<summary><h3>📋 Packages — 4 commands</h3></summary>

> Packages live inside offerings. They hold the products your paywall displays.

```bash
rc packages list --offering-id ofrngs_xxx               # List
rc packages create --offering-id ofrngs_xxx \
  --lookup-key monthly --display-name "Monthly"         # Create
rc packages attach-products \
  --package-id pkg_xxx --product-ids prod1,prod2        # Wire products
rc packages detach-products \
  --package-id pkg_xxx --product-ids prod1              # Unwire
```

</details>

<details>
<summary><h3>👤 Customers — 2 commands</h3></summary>

```bash
rc customers get --customer-id $APP_USER_ID             # Full customer info
rc customers delete --customer-id $APP_USER_ID --confirm
```

</details>

<details>
<summary><h3>🎨 Paywalls — 1 command</h3></summary>

```bash
rc paywalls create --offering-id ofrngs_xxx
```

</details>

<details>
<summary><h3>🔧 Utilities — 4 commands</h3></summary>

```bash
rc doctor                                               # Diagnose issues
rc init --project proj_xxx                              # Project config (.rc.yaml)
rc completion zsh > "${fpath[1]}/_rc"                   # Shell completions
rc version                                              # Version info
```

</details>

<br>

---

<br>

## 🖥️ Output Formats

Every command speaks 6 formats. Pick what fits.

| Flag | Format | Best for |
|:-----|:-------|:---------|
| *(default)* | JSON | Scripting, `jq`, APIs |
| `--pretty` | Pretty JSON | Reading in terminal |
| `-o table` | Aligned columns | Quick scanning |
| `-o csv` | CSV | Spreadsheets, exports |
| `-o tsv` | Tab-separated | Unix pipelines |
| `-o yaml` | YAML | Config files |
| `-o minimal` | First field only | Piping IDs |

```bash
# Same data, different shapes
rc entitlements list                  # [{"id":"entl_xxx","lookup_key":"premium",...}]
rc entitlements list -o table         # ID            LOOKUP_KEY    DISPLAY_NAME
rc entitlements list -o minimal       # entl_xxx
rc entitlements list -o csv           # ID,LOOKUP_KEY,DISPLAY_NAME
```

<br>

---

<br>

## 🔗 Scripting

`rc` is designed to be piped, parsed, and composed.

```bash
# Export your entire product catalog
rc products list --all -o csv > products.csv

# Count entitlements
rc entitlements list --all -o json | jq length

# Check if an offering exists
rc offerings list --all -o minimal | grep -q "default" && echo "exists"

# Get all lookup keys
rc entitlements list --all -o json | jq -r '.[].lookup_key'

# Paginate manually
rc products list --limit 10
rc products list --limit 10 --starting-after prod_xxx

# Or just get everything
rc products list --all
```

<br>

---

<br>

## ⚙️ Configuration

### Multiple Profiles

Switch between production, staging, or different projects seamlessly.

```bash
rc auth login --api-key sk_live_xxx --name production
rc auth login --api-key sk_test_xxx --name staging --default-project proj_staging
rc auth switch --name staging
rc auth list -o table
rc auth current
```

### Project Config (`.rc.yaml`)

Drop this in your repo root — `rc` auto-discovers it:

```bash
rc init --project proj_xxx --output table
```

```yaml
# .rc.yaml
project: proj_xxx
output: table
```

### Environment Variables

| Variable | Description |
|:---------|:------------|
| `RC_API_KEY` | API v2 secret key — overrides profile |
| `RC_PROJECT` | Project ID |
| `RC_PROFILE` | Active auth profile |
| `RC_OUTPUT` | Default output format |
| `RC_DEBUG` | `true` to log HTTP requests |
| `RC_TIMEOUT` | Request timeout (default `60s`) |

### Priority

Flags → Environment variables → `.rc.yaml` → `~/.revenuecat-cli/config.json`

<br>

---

<br>

## 🔒 Security

- Config stored at `~/.revenuecat-cli/config.json` with **`0600` permissions**
- API keys **masked** in `rc doctor` output
- Debug mode **redacts** Authorization headers
- In CI: use `RC_API_KEY` env var — nothing touches disk

<br>

---

<br>

## 🏗️ Built With

<table>
<tr>
<td align="center" width="25%"><a href="https://github.com/spf13/cobra"><strong>Cobra</strong></a><br><sub>CLI framework</sub></td>
<td align="center" width="25%"><a href="https://github.com/spf13/viper"><strong>Viper</strong></a><br><sub>Config management</sub></td>
<td align="center" width="25%"><a href="https://goreleaser.com"><strong>GoReleaser</strong></a><br><sub>Cross-platform releases</sub></td>
<td align="center" width="25%"><a href="https://www.revenuecat.com/docs/api-v2"><strong>RevenueCat API v2</strong></a><br><sub>REST API</sub></td>
</tr>
</table>

> Architecture follows [playconsole-cli](https://github.com/AndroidPoet/playconsole-cli) — the same patterns, proven at scale with 80+ commands.

<br>

---

<br>

## 🤝 Contributing

```bash
make build    # Build
make test     # Test
make lint     # Lint
make help     # All targets
```

PRs welcome. Open an issue first for major changes.

<br>

---

<br>

<div align="center">

**MIT License**

<sub>Not affiliated with RevenueCat Inc. RevenueCat is a trademark of RevenueCat Inc.</sub>

<br>
<br>

**[⭐ Star this repo](https://github.com/AndroidPoet/revenuecat-cli/stargazers)** if `rc` saved you a trip to the dashboard.

<br>

</div>

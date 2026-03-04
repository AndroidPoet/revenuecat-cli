# Apps

Manage applications within your RevenueCat project. Each app corresponds to a platform-specific distribution (App Store, Play Store, Stripe, etc.).

## Available Commands

| Command | Description |
|---|---|
| `rc apps list` | List all apps in the project |
| `rc apps get` | Get details for a specific app |
| `rc apps create` | Create a new app |
| `rc apps update` | Update an existing app |
| `rc apps delete` | Delete an app |
| `rc apps api-keys` | View API keys for an app |

---

## `rc apps list`

```bash
rc apps list
```

### Optional Flags

| Flag | Description | Default |
|---|---|---|
| `--output` | Output format | `json` |
| `--limit` | Maximum number of results | `20` |
| `--starting-after` | Cursor for pagination | -- |
| `--all` | Fetch all pages automatically | `false` |

---

## `rc apps get`

Retrieve details for a single app.

```bash
rc apps get --app-id app_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--app-id` | The app ID to retrieve |

---

## `rc apps create`

Create a new app in the current project.

```bash
rc apps create --name "My iOS App" --type app_store --bundle-id com.example.myapp
```

### Required Flags

| Flag | Description |
|---|---|
| `--name` | Display name for the app |
| `--type` | Store type |

### App Types

| Type | Description |
|---|---|
| `app_store` | Apple App Store |
| `play_store` | Google Play Store |
| `stripe` | Stripe |
| `amazon` | Amazon Appstore |
| `mac_app_store` | Mac App Store |
| `roku` | Roku Channel Store |
| `web` | Web |

### Optional Flags

| Flag | Description |
|---|---|
| `--bundle-id` | Bundle identifier (iOS/macOS) |
| `--package-name` | Package name (Android) |
| `--output` | Output format |

### Examples

=== "iOS App"

    ```bash
    rc apps create \
      --name "My iOS App" \
      --type app_store \
      --bundle-id com.example.myapp
    ```

=== "Android App"

    ```bash
    rc apps create \
      --name "My Android App" \
      --type play_store \
      --package-name com.example.myapp
    ```

=== "Stripe"

    ```bash
    rc apps create \
      --name "Web Payments" \
      --type stripe
    ```

---

## `rc apps update`

Update an existing app's properties.

```bash
rc apps update --app-id app_xxxxx --name "Renamed App"
```

### Required Flags

| Flag | Description |
|---|---|
| `--app-id` | The app ID to update |

### Optional Flags

| Flag | Description |
|---|---|
| `--name` | New display name |
| `--output` | Output format |

---

## `rc apps delete`

Delete an app from the project.

```bash
rc apps delete --app-id app_xxxxx --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--app-id` | The app ID to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    Deleting an app is permanent. All associated data will be lost.

---

## `rc apps api-keys`

View the public and secret API keys for an app.

```bash
rc apps api-keys --app-id app_xxxxx
```

### Required Flags

| Flag | Description |
|---|---|
| `--app-id` | The app ID |

!!! tip
    Use `--output minimal` to get just the key values, useful for piping into scripts or CI/CD workflows.

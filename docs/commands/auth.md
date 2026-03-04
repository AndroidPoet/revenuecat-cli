# Auth

Manage authentication profiles. Profiles store API keys locally and let you switch between different RevenueCat accounts or environments without re-entering credentials.

## Available Commands

| Command | Description |
|---|---|
| `rc auth login` | Create a new authentication profile |
| `rc auth switch` | Switch to a different profile |
| `rc auth list` | List all saved profiles |
| `rc auth current` | Show the currently active profile |
| `rc auth delete` | Delete a saved profile |

---

## `rc auth login`

Create and activate a new authentication profile.

```bash
rc auth login --api-key sk_xxxxx --name production
```

### Required Flags

| Flag | Description |
|---|---|
| `--api-key` | Your RevenueCat secret API key (starts with `sk_`) |
| `--name` | A name for this profile |

### Optional Flags

| Flag | Description |
|---|---|
| `--default-project` | Set a default project ID for this profile |

### Examples

=== "Basic"

    ```bash
    rc auth login --api-key sk_xxxxx --name production
    ```

=== "With Default Project"

    ```bash
    rc auth login \
      --api-key sk_xxxxx \
      --name production \
      --default-project proj_xxxxx
    ```

!!! note
    API keys are stored in your system's credential store. They are never written to plain-text config files.

---

## `rc auth switch`

Switch the active profile.

```bash
rc auth switch --name staging
```

### Required Flags

| Flag | Description |
|---|---|
| `--name` | The profile name to switch to |

---

## `rc auth list`

List all saved authentication profiles.

```bash
rc auth list
```

### Example Output

```
NAME          PROJECT        ACTIVE
production    proj_xxxxx     *
staging       proj_yyyyy
personal      --
```

---

## `rc auth current`

Display the currently active profile and its configuration.

```bash
rc auth current
```

### Example Output

```
Profile:   production
Project:   proj_xxxxx
API Key:   sk_...xxxxx (masked)
```

---

## `rc auth delete`

Delete a saved authentication profile.

```bash
rc auth delete --name old-profile --confirm
```

### Required Flags

| Flag | Description |
|---|---|
| `--name` | The profile name to delete |
| `--confirm` | Skip confirmation prompt |

!!! warning
    This permanently removes the profile and its stored API key. If this is the active profile, you will need to switch to another profile or create a new one.

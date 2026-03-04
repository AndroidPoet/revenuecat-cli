# Shell Completion

RevenueCat CLI provides tab completion for commands, subcommands, and flags. Set it up once and never type a full command name again.

## Setup

=== "Bash"

    Add to `~/.bashrc`:

    ```bash
    eval "$(rc completion bash)"
    ```

    Then reload your shell:

    ```bash
    source ~/.bashrc
    ```

=== "Zsh"

    Add to `~/.zshrc`:

    ```bash
    eval "$(rc completion zsh)"
    ```

    Then reload your shell:

    ```bash
    source ~/.zshrc
    ```

    !!! note
        If you get `command not found: compdef`, add this **before** the eval line:

        ```bash
        autoload -Uz compinit && compinit
        ```

=== "Fish"

    ```bash
    rc completion fish | source
    ```

    To make it permanent:

    ```bash
    rc completion fish > ~/.config/fish/completions/rc.fish
    ```

=== "PowerShell"

    Add to your PowerShell profile:

    ```powershell
    rc completion powershell | Out-String | Invoke-Expression
    ```

    To find your profile path:

    ```powershell
    echo $PROFILE
    ```

---

## What Gets Completed

Shell completion covers:

- **Commands**: `rc app` + Tab completes to `rc apps`
- **Subcommands**: `rc apps` + Tab shows `list`, `get`, `create`, `update`, `delete`, `api-keys`
- **Flags**: `rc apps create --` + Tab shows `--name`, `--type`, `--bundle-id`, etc.

---

## Verification

After setup, test that completion is working:

```bash
rc <Tab><Tab>
```

You should see a list of all available command groups:

```
apps           auth           audit-logs     customers
entitlements   init           metrics        offerings
packages       paywalls       products       projects
purchases      subscriptions  webhooks       doctor
version        completion
```

# Installation

## Homebrew (Recommended)

The fastest way to install RevenueCat CLI on macOS or Linux:

```bash
brew tap AndroidPoet/tap
brew install revenuecat-cli
```

To upgrade to the latest version:

```bash
brew upgrade revenuecat-cli
```

## Direct Download

Download pre-built binaries from the [GitHub Releases](https://github.com/AndroidPoet/revenuecat-cli/releases) page.

=== "macOS (Apple Silicon)"

    ```bash
    curl -LO https://github.com/AndroidPoet/revenuecat-cli/releases/latest/download/revenuecat-cli_darwin_arm64.tar.gz
    tar -xzf revenuecat-cli_darwin_arm64.tar.gz
    sudo mv rc /usr/local/bin/
    ```

=== "macOS (Intel)"

    ```bash
    curl -LO https://github.com/AndroidPoet/revenuecat-cli/releases/latest/download/revenuecat-cli_darwin_amd64.tar.gz
    tar -xzf revenuecat-cli_darwin_amd64.tar.gz
    sudo mv rc /usr/local/bin/
    ```

=== "Linux (amd64)"

    ```bash
    curl -LO https://github.com/AndroidPoet/revenuecat-cli/releases/latest/download/revenuecat-cli_linux_amd64.tar.gz
    tar -xzf revenuecat-cli_linux_amd64.tar.gz
    sudo mv rc /usr/local/bin/
    ```

=== "Linux (arm64)"

    ```bash
    curl -LO https://github.com/AndroidPoet/revenuecat-cli/releases/latest/download/revenuecat-cli_linux_arm64.tar.gz
    tar -xzf revenuecat-cli_linux_arm64.tar.gz
    sudo mv rc /usr/local/bin/
    ```

## Verify Installation

```bash
rc version
```

You should see output like:

```
revenuecat-cli version 0.1.0
```

## Next Steps

Once installed, proceed to [Configuration](configuration.md) to set up your authentication profile.

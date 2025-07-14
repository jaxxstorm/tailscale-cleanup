# Tailscale Cleanup Tool

> **⚠️ DEPRECATED**: This tool has been deprecated in favor of `tscli delete devices`. Please use [tscli](https://github.com/jaxxstorm/tscli) instead

A command-line utility for managing and cleaning up disconnected devices in your Tailscale tailnet. It identifies devices that have not been seen for a specified duration and removes them, with options for excluding certain devices and running in a dry-run mode for testing.

## Features

- List and identify disconnected devices based on the last seen duration.
- Exclude specific devices using partial name matching.
- Delete disconnected devices from your tailnet.
- Support for dry-run mode to preview changes without execution.
- Configurable via CLI flags or environment variables.

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Build the binary:
   ```bash
   go build -o tailscale-cleanup
   ```

3. Move the binary to a location in your `PATH` (optional):
   ```bash
   mv tailscale-cleanup /usr/local/bin/
   ```

## Usage

Run the tool with the required options:

```bash
tailscale-cleanup --api-key=<API_KEY> --tailnet=<TAILNET_NAME> [options]
```

### Required Flags

- `--api-key` or `TAILSCALE_API_KEY` (environment variable): Your Tailscale API key.
- `--tailnet` or `TAILNET_NAME` (environment variable): The name of your Tailscale tailnet.

### Optional Flags

- `--base-url`: Tailscale API base URL (default: `https://api.tailscale.com/api/v2`).
- `--last-seen-duration`: Duration to consider a device as disconnected (e.g., `15m`, `1h`, `24h`). Default is `15m`.
- `--exclude`: Device names or substrings to exclude from deletion. Can be specified multiple times.
- `--dry-run`: Run the tool without making destructive changes (default: `false`).

## Examples

### Basic Cleanup

```bash
tailscale-cleanup --api-key="your_api_key" --tailnet="example.com"
```

### Exclude Specific Devices

```bash
tailscale-cleanup --api-key="your_api_key" --tailnet="example.com" --exclude="server" --exclude="important-device"
```

### Dry Run

```bash
tailscale-cleanup --api-key="your_api_key" --tailnet="example.com" --dry-run
```

### Change Last Seen Duration

```bash
tailscale-cleanup --api-key="your_api_key" --tailnet="example.com" --last-seen-duration="24h"
```

## Environment Variables

The tool supports the following environment variables:

- `TAILSCALE_API_KEY`: Set your Tailscale API key.
- `TAILNET_NAME`: Set your Tailscale tailnet name.

## Output

The tool provides detailed output for each step, including:

- Skipped devices (due to exclusion or recent activity).
- Disconnected devices identified.
- Confirmation of deletions (or simulated actions in dry-run mode).

## Error Handling

If the tool encounters errors, they will be displayed in the output. Common errors include:

- Invalid API key or permissions.
- Network issues when connecting to the Tailscale API.
- Incorrectly formatted flags or environment variables.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with suggestions or improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


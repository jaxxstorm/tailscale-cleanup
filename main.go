package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
)

var Version = "dev"

// Device represents a Tailscale device
type Device struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"lastSeen"`
}

// ListResponse represents the API response for listing devices
type ListResponse struct {
	Devices []Device `json:"devices"`
}

// Config holds the CLI configuration
type Config struct {
	APIKey          string
	BaseURL         string
	TailnetName     string
	LastSeenTimeout time.Duration
	ExcludedDevices []string // List of substrings for partial match exclusion
}

func main() {
	var (
		app              = kingpin.New("tailscale-cleanup", "A utility to clean up disconnected Tailscale devices.")
		apiKey           = app.Flag("api-key", "Tailscale API key").Required().Envar("TAILSCALE_API_KEY").String()
		baseURL          = app.Flag("base-url", "Tailscale API base URL").Default("https://api.tailscale.com/api/v2").String()
		tailnetName      = app.Flag("tailnet", "Tailscale tailnet name").Required().Envar("TAILNET_NAME").String()
		lastSeenDuration = app.Flag("last-seen-duration", "Duration to consider a device disconnected (e.g., 15m, 1h)").Default("15m").Duration()
		exclude          = app.Flag("exclude", "Device names to exclude by partial match (can be specified multiple times)").Strings()
		dryRun           = app.Flag("dry-run", "Run without making destructive changes").Bool()
		showVersion      = kingpin.Flag("version", "Show the version and exit").Bool()
	)

	// Parse the command-line arguments
	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if *showVersion {
		fmt.Println("Version:", Version)
		os.Exit(0)
	}

	config := Config{
		APIKey:          *apiKey,
		BaseURL:         *baseURL,
		TailnetName:     *tailnetName,
		LastSeenTimeout: *lastSeenDuration,
		ExcludedDevices: *exclude,
	}

	err := cleanDisconnectedDevices(config, *dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func cleanDisconnectedDevices(config Config, dryRun bool) error {
	devices, err := listDevices(config)
	if err != nil {
		return fmt.Errorf("failed to list devices: %w", err)
	}

	now := time.Now()
	for _, device := range devices {
		// Check if the device is in the exclusion list by partial match
		if isExcluded(device.Name, config.ExcludedDevices) {
			fmt.Printf("Skipping excluded device: %s (%s)\n", device.Name, device.ID)
			continue
		}

		// Calculate time since the device was last seen
		timeSinceLastSeen := now.Sub(device.LastSeen)
		if timeSinceLastSeen > config.LastSeenTimeout {
			fmt.Printf("Disconnected device found (last seen %v ago): %s (%s)\n", timeSinceLastSeen, device.Name, device.ID)
			if !dryRun {
				err := deleteDevice(config, device.ID)
				if err != nil {
					fmt.Printf("Failed to delete device %s: %v\n", device.Name, err)
				} else {
					fmt.Printf("Deleted device: %s\n", device.Name)
				}
			} else {
				fmt.Println("Dry run enabled; skipping deletion.")
			}
		}
	}
	return nil
}

func isExcluded(deviceName string, excludedList []string) bool {
	for _, exclude := range excludedList {
		if strings.Contains(deviceName, exclude) {
			return true
		}
	}
	return false
}

func listDevices(config Config) ([]Device, error) {
	url := fmt.Sprintf("%s/tailnet/%s/devices", config.BaseURL, config.TailnetName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(config.APIKey, "")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %s, body: %s", resp.Status, body)
	}

	var listResponse ListResponse
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResponse.Devices, nil
}

func deleteDevice(config Config, deviceID string) error {
	url := fmt.Sprintf("%s/device/%s", config.BaseURL, deviceID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(config.APIKey, "")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %s, body: %s", resp.Status, body)
	}

	return nil
}

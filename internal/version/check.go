package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	GitHubAPIURL = "https://api.github.com/repos/atharvamhaske/hashctl/releases/latest"
	CheckTimeout = 5 * time.Second
)

// ReleaseInfo contains GitHub release information
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	URL     string `json:"html_url"`
}

// CheckLatestVersion checks GitHub for the latest release version
func CheckLatestVersion(currentVersion string) (*ReleaseInfo, error) {
	client := &http.Client{
		Timeout: CheckTimeout,
	}

	req, err := http.NewRequest("GET", GitHubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

// IsUpdateAvailable checks if a newer version is available
func IsUpdateAvailable(currentVersion, latestVersion string) bool {
	// Remove 'v' prefix for comparison
	current := strings.TrimPrefix(currentVersion, "v")
	latest := strings.TrimPrefix(latestVersion, "v")

	// Simple comparison - if versions are different, assume update available
	return current != latest
}

// GetUpdateMessage returns a formatted update notification message
func GetUpdateMessage(currentVersion, latestVersion, releaseURL string) string {
	return fmt.Sprintf(
		" Update available: %s → %s\n   Download: %s",
		currentVersion,
		latestVersion,
		releaseURL,
	)
}

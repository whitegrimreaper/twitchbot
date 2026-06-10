// Package twitchtoken handles refreshing Twitch OAuth tokens and persisting
// them to a .env file.
package twitchtoken

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	// Populated on error
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// GetAppToken requests a fresh app access token using the client credentials
// flow. This is used for server-to-server API calls (Helix/EventSub) and does
// not require a refresh token — just the app's client ID and secret.
func GetAppToken(clientID, clientSecret string) (RefreshResult, error) {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return RefreshResult{}, fmt.Errorf("app token HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var result refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RefreshResult{}, fmt.Errorf("failed to parse app token response: %w", err)
	}
	if result.AccessToken == "" {
		return RefreshResult{}, fmt.Errorf("twitch app token request failed (status %d) with error: %s", result.Status, result.Message)
	}

	return RefreshResult{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
		// Client credentials tokens have no refresh token
	}, nil
}

// Refresh exchanges a refresh token for a new access token using the
// Twitch authorization code flow. Used for user-scoped tokens (e.g. IRC).
func Refresh(clientID, clientSecret, refreshToken string) (RefreshResult, error) {
	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	})
	if err != nil {
		return RefreshResult{}, fmt.Errorf("token refresh HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var result refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RefreshResult{}, fmt.Errorf("failed to parse token refresh response: %w", err)
	}
	if result.AccessToken == "" {
		return RefreshResult{}, fmt.Errorf("twitch token refresh failed (status %d) with error: %s", result.Status, result.Message)
	}

	return RefreshResult{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// UpdateEnvFile rewrites the .env file at path, replacing values for the
// given keys. Keys not already present in the file are appended at the end.
func UpdateEnvFile(path string, updates map[string]string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var lines []string
	updated := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for key, val := range updates {
			if strings.HasPrefix(line, key+"=") {
				line = key + "=" + val
				updated[key] = true
				break
			}
		}
		lines = append(lines, line)
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	for key, val := range updates {
		if !updated[key] {
			lines = append(lines, key+"="+val)
		}
	}

	var buf bytes.Buffer
	for _, line := range lines {
		buf.WriteString(line + "\n")
	}

	return os.WriteFile(path, buf.Bytes(), 0600)
}

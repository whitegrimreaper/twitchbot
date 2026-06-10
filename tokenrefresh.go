package main

import (
	"fmt"
	"os"

	"bot/internal/twitchtoken"
)

// refreshAppToken requests a fresh Helix app access token using the client
// credentials flow, then updates both the running process env and the .env file.
// Returns the new raw token string (no "oauth:" prefix).
func refreshAppToken(envPath string) (string, error) {
	clientID := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("APP_CLIENT_ID and APP_CLIENT_SECRET must be set in .env")
	}

	fmt.Println("Attempting to refresh Helix app access token...")

	result, err := twitchtoken.GetAppToken(clientID, clientSecret)
	if err != nil {
		return "", err
	}

	os.Setenv("ACCESS_TOKEN_FINAL", result.AccessToken)

	if err := twitchtoken.UpdateEnvFile(envPath, map[string]string{
		"ACCESS_TOKEN_FINAL": result.AccessToken,
	}); err != nil {
		fmt.Printf("Warning: app token refreshed in memory but failed to write to %s: %v\n", envPath, err)
		fmt.Printf("New ACCESS_TOKEN_FINAL: %s\n", result.AccessToken)
	} else {
		fmt.Printf("App token refreshed successfully. Expires in ~%.1f days.\n", float64(result.ExpiresIn)/86400)
	}

	return result.AccessToken, nil
}

// refreshIRCToken uses the stored refresh token to get a new IRC access token
// from Twitch, then updates both the running process env and the .env file on disk.
// Returns the new token in "oauth:xxxx" format ready for use with the IRC client.
func refreshIRCToken(envPath string) (string, error) {
	clientID := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")
	refreshToken := os.Getenv("AUTH_TOKEN_REFRESH")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("APP_CLIENT_ID and APP_CLIENT_SECRET must be set in .env")
	}
	if refreshToken == "" {
		return "", fmt.Errorf("AUTH_TOKEN_REFRESH is not set in .env — cannot auto-refresh")
	}

	fmt.Println("Attempting to refresh IRC token...")

	result, err := twitchtoken.Refresh(clientID, clientSecret, refreshToken)
	if err != nil {
		return "", err
	}

	newAuthToken := "oauth:" + result.AccessToken

	// Update the running process so anything else reading these env vars is current
	os.Setenv("AUTH_TOKEN", newAuthToken)
	os.Setenv("AUTH_TOKEN_REFRESH", result.RefreshToken)

	// Persist to disk so the next run picks up the new tokens
	if err := twitchtoken.UpdateEnvFile(envPath, map[string]string{
		"AUTH_TOKEN":         newAuthToken,
		"AUTH_TOKEN_REFRESH": result.RefreshToken,
	}); err != nil {
		// Not fatal — the in-memory tokens are already updated and the bot can continue.
		fmt.Printf("Warning: tokens refreshed in memory but failed to write to %s: %v\n", envPath, err)
		fmt.Printf("New AUTH_TOKEN: %s\n", newAuthToken)
		fmt.Printf("New AUTH_TOKEN_REFRESH: %s\n", result.RefreshToken)
	} else {
		fmt.Printf("Token refreshed successfully. Expires in ~%.1f hours.\n", float64(result.ExpiresIn)/3600)
	}

	return newAuthToken, nil
}

package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"bot/internal/twitchtoken"
)

func main() {
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("No .env file found\n")
		os.Exit(1)
	}

	clientID := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")
	refreshToken := os.Getenv("AUTH_TOKEN_REFRESH")

	validateEnvVars([]string{clientID, clientSecret, refreshToken})

	fmt.Println("Requesting new token from Twitch...")

	result, err := twitchtoken.Refresh(clientID, clientSecret, refreshToken)
	if err != nil {
		fmt.Printf("Token refresh failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("New access token:  %s\n", result.AccessToken)
	fmt.Printf("New refresh token: %s\n", result.RefreshToken)
	fmt.Printf("Expires in:        %d seconds (~%.1f hours)\n", result.ExpiresIn, float64(result.ExpiresIn)/3600)

	if err := twitchtoken.UpdateEnvFile(envPath, map[string]string{
		"AUTH_TOKEN":         "oauth:" + result.AccessToken,
		"AUTH_TOKEN_REFRESH": result.RefreshToken,
	}); err != nil {
		fmt.Printf("Tokens fetched but failed to update .env: %v\n", err)
		fmt.Println("Update manually using the values printed above.")
		os.Exit(1)
	}

	fmt.Println("\n.env updated successfully.")
}

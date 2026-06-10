package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"bot/internal/twitchtoken"
)

func main() {
	envPath := ".env"

	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Could not load .env file: %v\n", err)
		os.Exit(1)
	}

	clientID := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		fmt.Println("APP_CLIENT_ID and APP_CLIENT_SECRET must be set in .env")
		os.Exit(1)
	}

	fmt.Println("Requesting new app access token from Twitch...")

	result, err := twitchtoken.GetAppToken(clientID, clientSecret)
	if err != nil {
		fmt.Printf("App token refresh failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("New app access token: %s\n", result.AccessToken)
	fmt.Printf("Expires in:           %d seconds (~%.1f days)\n", result.ExpiresIn, float64(result.ExpiresIn)/86400)

	if err := twitchtoken.UpdateEnvFile(envPath, map[string]string{
		"ACCESS_TOKEN_FINAL": result.AccessToken,
	}); err != nil {
		fmt.Printf("Token fetched but failed to update .env: %v\n", err)
		fmt.Println("Update ACCESS_TOKEN_FINAL manually using the value printed above.")
		os.Exit(1)
	}

	fmt.Println("\n.env updated successfully.")
}

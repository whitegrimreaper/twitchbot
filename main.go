package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found\n")
	}
}

func main() {
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		fmt.Printf("AUTH_TOKEN not set\n")
		os.Exit(0)
	}
	fmt.Printf("AUTH_TOKEN is set to %s\n", authToken)

	// Twitch bot configuration
	botUsername := "whitescancerbot"
	channel := "whitegrimreaper_"

	// Create a new Twitch IRC client
	fmt.Printf("INITIALIZING BOT\n")
	client := twitch.NewClient(botUsername, authToken)
	client.Join(channel)

	// Event handler for incoming messages
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("[%s] %s: %s\n", message.Channel, message.User.DisplayName, message.Message)

		isCommand, err := isCommand(message)
		if err == nil && isCommand {
			fmt.Printf("Found something matching command syntax!\n")
			args, err := extractArgs(message)
			if err == nil {
				fmt.Printf("Extracted command %s\n", args[0])
				response, err2 := handleCommand(args, message)
				fmt.Printf("Should respond: '%s'\n", response)
				if err2 == nil && response != "" {
					client.Say(channel, response)
				}
			}
		}
	})

	// Connect to Twitch IRC
	err := client.Connect()
	if err != nil {
		fmt.Printf("Error connecting to Twitch IRC: %v\n", err)
		return
	}
	fmt.Printf("Connected to Twitch IRC\n")

	// Wait for termination signal (Ctrl+C)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Disconnect from Twitch IRC
	client.Disconnect()
	fmt.Printf("Disconnected from Twitch IRC\n")
}

func extractArgs(message twitch.PrivateMessage) (args []string, err error) {
	ret := strings.Fields(message.Message)
	//for i := 0; i < len(ret); i++ {
	//	fmt.Printf("%s,", ret[i])
	//}
	//fmt.Printf("\n")
	return ret, nil
}

func isCommand(message twitch.PrivateMessage) (ret bool, err error) {
	parsedMessage := message.Message
	if parsedMessage[0] == markerChar {
		return true, nil
	}
	return false, nil
}

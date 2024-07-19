package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix/v2"
	"gorm.io/gorm"
	//"gorm.io/driver/sqlite"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

var BossDB *gorm.DB
var PointsDB *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found\n")
	}
}

func main() {
	authToken := os.Getenv("AUTH_TOKEN")
	clientId := os.Getenv("APP_CLIENT_ID")
	clientSecret := os.Getenv("APP_CLIENT_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN_FINAL")
	validateEnvVars([]string{authToken, clientId, clientSecret, accessToken})
	fmt.Printf("AUTH_TOKEN is set to %s\n", authToken)

	PointsDB = pointsDBInit()
	BossDB = bossDBInit()

	helixClient, err := helixInit(clientId, clientSecret, accessToken)
	if err != nil {
		panic(err)
	}
	createEventSubSubscriptions(helixClient)

	resp, err := helixClient.GetUsers(&helix.UsersParams{
		Logins: []string{"whitegrimreaper_"},
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Message: %s\n", resp.ErrorMessage)
	fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())
	
	for _, user := range resp.Data.Users {
		fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
	}

	// TODO currently I create new listeners every time the bot is run, so I'm at 13 subscriptions as of the writing of this
	// should really remove them here

	//go startTwitchEventListener()
	go startTwitchListeners()

	time.Sleep(3*time.Second)
	// leaving this in as a reminder to check eventsub subs at some point
	/*eventSubResp, err := helixClient.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		Status: helix.EventSubStatusEnabled, // This is optional.
	})
	if err != nil {
		panic(err)
	}*/

	// Twitch bot configuration
	botUsername := "whitegrimbot"
	channel := "whitegrimreaper_"

	// Create a new Twitch IRC client
	fmt.Printf("Initializing Bot\n")
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
	conErr := client.Connect()
	// always hangs here, never gets to the following code
	// i think client.Connect() should be called in an asynch manner from a
	// separate goroutine but that's a future improvement since we don't care atm
	if conErr != nil {
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
	return ret, nil
}

func isCommand(message twitch.PrivateMessage) (ret bool, err error) {
	parsedMessage := message.Message
	if parsedMessage[0] == markerChar {
		return true, nil
	}
	return false, nil
}

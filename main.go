package main

import (
	"context"
	"fmt"
	//"time"
	"os"
	"strings"
	"net/http"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"

	"github.com/joho/godotenv"
	"github.com/nicklaw5/helix/v2"
	"gorm.io/gorm"
	//"gorm.io/driver/sqlite"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

var BossDB *gorm.DB
var ReqQueueDB *gorm.DB
var BossNickDB *gorm.DB
var PointsDB *gorm.DB
var WebhookClient *http.Client

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

	region := "prod" // default to 'prod'

	// this makes it easy to fuck with testing stuff
	if len(os.Args) > 1 {
		if os.Args[1] == "test" {
			region = "test"
		}
	}

	PointsDB = pointsDBInit(region)
	ReqQueueDB = reqQueueDBInit(region)
	BossNickDB = bossNicksDBInit()
	BossDB = bossDBInit()
	WebhookClient = webhookClientInit()

	helixClient, err := helixInit(clientId, clientSecret, accessToken)
	if err != nil {
		panic(err)
	}

	// leaving this in as a reminder to check eventsub subs at some point, should really just update
	// existing ones instead of remaking every time
	eventSubResp, err := helixClient.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{
		//Status: helix.EventSubStatusEnabled, // This is optional.
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current number of eventsub subs: %d :)\n\n",eventSubResp.Data.Total)
	if eventSubResp.Data.Total > 4 {
		// if we have more than 4, delete em all and refresh
		for _, sub := range eventSubResp.Data.EventSubSubscriptions {
			deleteEventSubSubscription(helixClient, sub.ID)
		}
		createEventSubSubscriptions(helixClient)
	} else if eventSubResp.Data.Total < 4 {
		// if less than 4, same idea
		// if 1-3 already exist, a repeat call should just do nothing
		createEventSubSubscriptions(helixClient)
	}

	resp, err := helixClient.GetUsers(&helix.UsersParams{
		Logins: []string{"whitegrimreaper_"},
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("\n==Testing Helix Return==\n")
	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n", resp.GetRateLimitReset())
	
	for _, user := range resp.Data.Users {
		fmt.Printf("ID: %s Name: %s\n", user.ID, user.DisplayName)
	}
	fmt.Printf("==Done with Helix Return==\n\n")

	// TODO currently I create new listeners every time the bot is run, so I'm at 91 subscriptions (updating this every time i see)
	// should really remove them here
	go startTwitchListeners()

	// Listener For StreamDeck Events
	go startStreamDeckListener()

	// Listener for queue visual
	go startBossQueueListener()

	// TODO: Should move the overlay/browser stuff to this package
	// mostly so I don't have to run multiple separate go programs
	// to get things running

	// Twitch bot configuration
	botUsername := "whitegrimbot"
	channel := "whitegrimreaper_"

	// Connect to Twitch IRC, refresh token if auth fails
	fmt.Printf("Initializing Bot\n")
	conErr := connectIRC(botUsername, channel, authToken)
	if conErr != nil {
		fmt.Printf("Error connecting to Twitch IRC: %v\n", conErr)
		fmt.Printf("Attempting token refresh...\n")

		newToken, refreshErr := refreshIRCToken(".env")
		if refreshErr != nil {
			fmt.Printf("Token refresh failed: %v\n", refreshErr)
			return
		}

		fmt.Printf("Retrying connection with refreshed token...\n")
		conErr = connectIRC(botUsername, channel, newToken)
		if conErr != nil {
			fmt.Printf("Connection failed after token refresh: %v\n", conErr)
			return
		}
	}
	fmt.Printf("Disconnected from Twitch IRC\n")
}

// slightly better way of doing things
func connectIRC(botUsername, channel, token string) error {
	client := twitch.NewClient(botUsername, token)
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

	return client.Connect()
}

func extractArgs(message twitch.PrivateMessage) (args []string, err error) {
	// This is a really weird glitch of some kind, not sure if Twitch or go-twitch-irc
	// but sometimes (always with identical messages), messages will have this
	// Unknown character appended to the end. So we trim it here. Idk why it does that
	ret := strings.Fields(strings.Trim(message.Message, " 󠀀"))
	return ret, nil
}

func isCommand(message twitch.PrivateMessage) (ret bool, err error) {
	parsedMessage := message.Message
	if parsedMessage[0] == markerChar {
		return true, nil
	}
	return false, nil
}

func run(ctx context.Context) error {
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	fmt.Printf("Ingress established at: %s\n", listener.URL())

	return http.Serve(listener, http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from ngrok-go!")
}
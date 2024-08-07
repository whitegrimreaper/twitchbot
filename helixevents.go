package main

import (
	"fmt"
	"strconv"
	"os"
	"net/http"
	"encoding/json"
	"bytes"
	"io"
	helix "github.com/nicklaw5/helix/v2"
)

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func helixInit(clientId, clientSecret, accessToken string)(client *helix.Client, err error) {
	helixClient, err := helix.NewClient(&helix.Options{
		ClientID: clientId,
		ClientSecret: clientSecret,
		AppAccessToken: accessToken,
	})
	if err != nil {
		panic(err)
	}
	return helixClient, err
}

func startTwitchListeners() {
	http.HandleFunc("/twitch/webhook", handleTwitchCallback)
	http.HandleFunc("/eventsub", eventSubHandler)
	http.ListenAndServe(":3000", nil)
}

func handleTwitchCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Process the incoming webhook payload
		// Verify the signature if Twitch sends one
		// Handle the event data
		fmt.Println("Received Twitch webhook:", r.Body)
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func createEventSubSubscriptions(client *helix.Client) {
	fmt.Printf("\n=====================================================\nCREATING EVENTSUB SUBSCRIPTIONS\n======================================================\n");
	fmt.Printf("Creating follow subscription...\n")
	resp, err := client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type: helix.EventSubTypeChannelFollow,
		Version: "2",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: os.Getenv("BROADCASTER_ID"),
			ModeratorUserID: os.Getenv("BOT_ID"),
		},
		Transport: helix.EventSubTransport{
			Method: "webhook",
			Callback: "https://localhost:443",
			Secret: os.Getenv("EVENTSUB_SECRET"),
		},
	})
	if err != nil {
		panic(err)
	}
	printThingsFromEventSubResp(resp)

	fmt.Printf("Creating subscribe subscription...\n")
	resp, err = client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type: helix.EventSubTypeChannelSubscription,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: os.Getenv("BROADCASTER_ID"),
			ModeratorUserID: os.Getenv("BOT_ID"),
		},
		Transport: helix.EventSubTransport{
			Method: "webhook",
			Callback: "https://localhost:443",
			Secret: os.Getenv("EVENTSUB_SECRET"),
		},
	})
	if err != nil {
		panic(err)
	}
	printThingsFromEventSubResp(resp)

	fmt.Printf("Creating gift sub subscription...\n")
	resp, err = client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type: helix.EventSubTypeChannelSubscriptionGift,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: os.Getenv("BROADCASTER_ID"),
			ModeratorUserID: os.Getenv("BOT_ID"),
		},
		Transport: helix.EventSubTransport{
			Method: "webhook",
			Callback: "https://localhost:443",
			Secret: os.Getenv("EVENTSUB_SECRET"),
		},
	})
	if err != nil {
		panic(err)
	}
	printThingsFromEventSubResp(resp)

	fmt.Printf("=====================================================\nDONE WITH EVENTSUB SUBSCRIPTIONS\n======================================================\n");
}

// Used for deleting all ESS's because currently I don't handle them well
func deleteEventSubSubscription(client *helix.Client, subId string) {
	resp, err := client.RemoveEventSubSubscription(subId)
	if err != nil {
		fmt.Printf("Error removing EventSubSub %s : %s\n",subId, err)
	}
	fmt.Printf("ESS %+v removed\n", resp)
}

func eventSubHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("%s\n",err)
        return
    }
    defer r.Body.Close()
    // verify that the notification came from twitch using the secret.
    if !helix.VerifyEventSubNotification(os.Getenv("EVENTSUB_SECRET"), r.Header, string(body)) {
        fmt.Printf("no valid signature on subscription")
        return
    } else {
        fmt.Printf("verified signature for subscription\n")
    }
    var vals eventSubNotification
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
    if err != nil {
        fmt.Printf("%d\n",err)
        return
    }
	// this line prints the whole notif if we wanna see it, but we probably don't every time
	// fmt.Printf("\n\n\n%+v\n\n\n",vals)
    // if there's a challenge in the request, respond with only the challenge to verify your eventsub.
    if vals.Challenge != "" {
        w.Write([]byte(vals.Challenge))
        return
    }

	switch vals.Subscription.Type {
	case helix.EventSubTypeChannelFollow:
		var followEvent helix.EventSubChannelFollowEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&followEvent)
		if err != nil {fmt.Printf("")}
		fmt.Printf("Got Follow event!\n")
		fmt.Printf("%s follows %s\n", followEvent.UserName, followEvent.BroadcasterUserName)
	case helix.EventSubTypeChannelSubscription:
		var subEvent helix.EventSubChannelSubscribeEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&subEvent)
		if err != nil {fmt.Printf("")}
		fmt.Printf("Got Sub event!\n")
		fmt.Printf("%s subbed to %s, and was it gift? %t\n", subEvent.UserName, subEvent.BroadcasterUserName, subEvent.IsGift)
		idInt, err := strconv.Atoi(subEvent.UserID)
		if err != nil {
			panic(err)
		}
		respCode, respMessage, exists := doesUserExist(idInt)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("")
		}
		fmt.Printf("Does user exist? %t\n", exists)
		respCode, respMessage = writePointGainEvent(idInt, 100)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("")
		}
	case helix.EventSubTypeChannelCheer:
		var cheerEvent helix.EventSubChannelCheerEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&cheerEvent)
		if err != nil {fmt.Printf("")}
		fmt.Printf("Got Cheer event!\n")
		fmt.Printf("%s gave bits to %s in amount %d!\n", cheerEvent.UserName, cheerEvent.BroadcasterUserName, cheerEvent.Bits)
	}

    //fmt.Printf("got follow webhook: %s follows %s\n", followEvent.UserName, followEvent.BroadcasterUserName)
    w.WriteHeader(200)
    w.Write([]byte("ok"))
}
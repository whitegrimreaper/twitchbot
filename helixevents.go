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

type StreamDeckPayload struct {
	BossName string `json:"boss_name"`
	NumKills    int `json:"num_kills"`
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

func startStreamDeckListener() {
	http.HandleFunc("/webhook", streamDeckWebhookHandler)
	http.ListenAndServe(":8081", nil)
}

func startSecureTwitchListeners() {
	http.HandleFunc("/eventsub", eventSubHandler)
	fmt.Println("Starting server on port 443")
    err := http.ListenAndServeTLS(":443", "certs/cert.pem", "certs/key.pem", nil)
    if err != nil {
        panic(err)
    }
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
	fmt.Printf("\n================================================\nCREATING EVENTSUB SUBSCRIPTIONS\n================================================\n");
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
			Callback: "https://possum-subtle-quagga.ngrok-free.app/eventsub",
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
			Callback: "https://possum-subtle-quagga.ngrok-free.app/eventsub",
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
			Callback: "https://possum-subtle-quagga.ngrok-free.app/eventsub",
			Secret: os.Getenv("EVENTSUB_SECRET"),
		},
	})
	if err != nil {
		panic(err)
	}
	printThingsFromEventSubResp(resp)

	fmt.Printf("Creating channel point subscription...\n")
	resp, err = client.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type: helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: os.Getenv("BROADCASTER_ID"),
			ModeratorUserID: os.Getenv("BOT_ID"),
		},
		Transport: helix.EventSubTransport{
			Method: "webhook",
			Callback: "https://possum-subtle-quagga.ngrok-free.app/eventsub",
			Secret: os.Getenv("EVENTSUB_SECRET"),
		},
	})
	if err != nil {
		panic(err)
	}
	printThingsFromEventSubResp(resp)

	fmt.Printf("================================================\nDONE WITH EVENTSUB SUBSCRIPTIONS\n================================================\n");
}

// Used for deleting all ESS's if we need to refresh
func deleteEventSubSubscription(client *helix.Client, subId string) {
	resp, err := client.RemoveEventSubSubscription(subId)
	if err != nil {
		fmt.Printf("Error removing EventSubSub %s : %s\n",subId, err)
	}
	fmt.Printf("ESS %+v removed\n", resp)
}

func eventSubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got some sort of eventsub thing\n")
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
		fmt.Printf("%s followed the channel!\n", followEvent.UserName)
		idInt, err := strconv.Atoi(followEvent.UserID)
		if err != nil {
			panic(err)
		}
		respCode, respMessage := writePointGainEvent(idInt, followEvent.UserName,25)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("")
		}
		sendPointGainWebook(followEvent.UserName, 25, "following")
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
		respCode, respMessage := writePointGainEvent(idInt, subEvent.UserName,500)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("")
		}
		if subEvent.IsGift {
			// also give gifter points
			var giftEvent helix.EventSubChannelSubscriptionGiftEvent 
			err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&giftEvent)
			if !giftEvent.IsAnonymous {
					idInt, err := strconv.Atoi(subEvent.UserID)
				if err != nil {
					panic(err)
				}
				respCode, respMessage := writePointGainEvent(idInt, giftEvent.UserName, 500 * giftEvent.Total)
				if respCode != 0 || respMessage != "" {
					fmt.Printf("")
				}
			}
		}
		sendPointGainWebook(subEvent.UserName, 500, "subbing/gifting")
	case helix.EventSubTypeChannelCheer:
		var cheerEvent helix.EventSubChannelCheerEvent
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&cheerEvent)
		if err != nil {fmt.Printf("")}
		fmt.Printf("Got Cheer event!\n")
		fmt.Printf("%s gave %d bits!\n", cheerEvent.UserName, cheerEvent.Bits)
		idInt, err := strconv.Atoi(cheerEvent.UserID)
			if err != nil {
				panic(err)
			}
			respCode, respMessage := writePointGainEvent(idInt, cheerEvent.UserName, cheerEvent.Bits)
			if respCode != 0 || respMessage != "" {
				fmt.Printf("")
			}
			sendPointGainWebook(cheerEvent.UserName, cheerEvent.Bits, "cheering")
	case helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd:
		var pointsEvent helix.EventSubChannelPointsCustomRewardRedemptionEvent 
		err = json.NewDecoder(bytes.NewReader(vals.Event)).Decode(&pointsEvent)
		if err != nil {fmt.Printf("")}
		fmt.Printf("Got channel point event!\n")
		fmt.Printf("%s redeemed %d points for %s!\n", pointsEvent.UserName, pointsEvent.Reward.Cost, pointsEvent.Reward.Title)
		if(pointsEvent.Reward.Title == "+25 Boss Points") {
			idInt, err := strconv.Atoi(pointsEvent.UserID)
			if err != nil {
				panic(err)
			}
			respCode, respMessage := writePointGainEvent(idInt, pointsEvent.UserName, 25)
			if respCode != 0 || respMessage != "" {
				fmt.Printf("")
			}
			sendPointGainWebook(pointsEvent.UserName, 25, "redeeming points")
		} else if(pointsEvent.Reward.Title == "+100 Boss Points") {
			idInt, err := strconv.Atoi(pointsEvent.UserID)
			if err != nil {
				panic(err)
			}
			respCode, respMessage := writePointGainEvent(idInt, pointsEvent.UserName, 100)
			if respCode != 0 || respMessage != "" {
				fmt.Printf("")
			}
			sendPointGainWebook(pointsEvent.UserName, 100, "redeeming points")
		} else if(pointsEvent.Reward.Title == "+200 Boss Points") {
			idInt, err := strconv.Atoi(pointsEvent.UserID)
			if err != nil {
				panic(err)
			}
			respCode, respMessage := writePointGainEvent(idInt, pointsEvent.UserName, 200)
			if respCode != 0 || respMessage != "" {
				fmt.Printf("")
			}
			sendPointGainWebook(pointsEvent.UserName, 200, "redeeming points")
		}
	}

    //fmt.Printf("got follow webhook: %s follows %s\n", followEvent.UserName, followEvent.BroadcasterUserName)
    w.WriteHeader(200)
    w.Write([]byte("ok"))
}

func streamDeckWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload StreamDeckPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received payload: %+v\n", payload)
	resp := executeBossRemoval(payload.BossName, payload.NumKills)
	if resp != "Success!" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func startBossQueueListener() {
	http.HandleFunc("/oldest-requests", getOldestRequestsHandler())
    http.ListenAndServe(":8082", nil)
}



func getOldestRequestsHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var requests []UserBossRequest

        // Query the 5 oldest entries ordered by CreatedAt
        err := ReqQueueDB.Order("created_at ASC").Limit(5).Find(&requests).Error
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(requests)
    }
}
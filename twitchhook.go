package main
// This code is not verified yet, don't have callback enabled from twitch API
// but it's still here so i can deploy onto gcloud and get the callback url from there

// NOTE: PROBABLY GONNA TRASH THIS SINCE THE HELIX THING WORKS AND IS MUCH EASIER TO INTERACT WITH
// LATER NOTE: DEF TRASHING THIS BUT SAVING IT FOR LATER SO MY GIT HISTORY LOOKS NICER

import (
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

// Structs to unmarshal Twitch EventSub payloads

// Top-level structure for Twitch EventSub notification
type EventSubNotification struct {
	Subscription  Subscription       `json:"subscription"`
	Event         json.RawMessage    `json:"event"`
}

type TwitchEventSubscription struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Condition struct {
		BroadcasterUserId string `json:"broadcaster_user_id"`
	} `json:"condition"`
}

// Subscription details
type Subscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Version   string `json:"version"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
	} `json:"condition"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
		Secret   string `json:"secret"`
	} `json:"transport"`
	CreatedAt string `json:"created_at"`
}

func subscribeToFollowEvents() {
	// Example subscription payload for follows to a specific broadcaster (replace with your broadcaster id)
	subscription := TwitchEventSubscription{
		Type:    "channel.follow",
		Version: "1",
		Condition: struct {
			BroadcasterUserId string `json:"broadcaster_user_id"`
		}{
			BroadcasterUserId: "123456789", // Replace with your broadcaster id
		},
	}

	// Convert subscription to JSON
	payload, err := json.Marshal(subscription)
	if err != nil {
		log.Fatalf("Error marshalling subscription payload: %v", err)
	}

	// Make HTTP request to Twitch API to create subscription
	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %v", err)
	}
	req.Header.Set("Client-ID", "YOUR_TWITCH_CLIENT_ID") // Replace with your Twitch Client ID
	req.Header.Set("Authorization", "Bearer YOUR_TWITCH_ACCESS_TOKEN") // Replace with your Twitch OAuth token
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Subscription response: %v", resp.Status)
}

// Handle incoming Twitch EventSub notifications
func handleTwitchEventSubNotification(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify Twitch signature
	signature := r.Header.Get("Twitch-Eventsub-Message-Signature")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	if !isValidSignature(signature, timestamp, body) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	var notification EventSubNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Extract event type and handle accordingly
	switch notification.Subscription.Type {
	case "channel.subscribe":
		handleChannelSubscribe(notification)
	case "channel.follow":
		handleChannelFollow(notification)
	case "channel.bits":
		handleChannelBits(notification)
	default:
		log.Printf("Unhandled EventSub type: %s", notification.Subscription.Type)
	}

	w.WriteHeader(http.StatusOK)
}

// checks signature against wehbook secret
func isValidSignature(signature, timestamp string, body []byte) bool {
	// Compute HMAC SHA256 hash using the webhook secret and payload
	hash := hmac.New(sha256.New, []byte(os.Getenv("WEBHOOK_SECRET")))
	hash.Write([]byte(fmt.Sprintf("%s%s", timestamp, string(body))))
	expectedMAC := hex.EncodeToString(hash.Sum(nil))

	// Compare the computed HMAC with the provided signature
	return hmac.Equal([]byte(signature), []byte("sha256="+expectedMAC))
}

// just logging for now since it's not functional anyways
// Handle channel subscribe event
func handleChannelSubscribe(notification EventSubNotification) {
	log.Printf("New subscription event: %+v", notification)
}

// Handle channel follow event
func handleChannelFollow(notification EventSubNotification) {
	log.Printf("New follow event: %+v", notification)
}

// Handle channel bits event
func handleChannelBits(notification EventSubNotification) {
	log.Printf("New bits donation event: %+v", notification)
}

func startTwitchEventListener() {
	// TODO: Add callback listener
	mux := http.NewServeMux()
	mux.HandleFunc("/twitch/webhook", handleTwitchEventSubNotification)

	// needs some work
	port := ":8080"
	fmt.Printf("Starting server on %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	select{}
}
// This class holds webhook stuff

package main

import (
	"bytes"
	"fmt"
	"strconv"
	"encoding/json"
	"net/http"
)

const eventWebhookURL = "http://localhost:8080/webhook"

func webhookClientInit()(client *http.Client) {
	client = &http.Client{}
	return client
}

func sendPointRedemptionWebhook(user string, boss string, numKills int) {
	payload := map[string]interface{}{
		"message": user + " has added " + strconv.Itoa(numKills) + " " + boss + " kills!",
	}

	payloadBytes, err := json.Marshal(payload)
    if err != nil {
        fmt.Println("Error marshalling payload:", err)
        return
    }

    // Create a new HTTP POST request
    req, err := http.NewRequest("POST", eventWebhookURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    resp, err := WebhookClient.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

	fmt.Println("Response From Internal Webhook:", resp.Status)
}

func sendPointGainWebook(user string, numPoints int, reason string) {
	fmt.Printf("Sending point gain webhook!\n")
	payload := map[string]interface{}{
		"message": user + " gained " + strconv.Itoa(numPoints) + " points for " + reason + "!",
	}

	payloadBytes, err := json.Marshal(payload)
    if err != nil {
        fmt.Println("Error marshalling payload:", err)
        return
    }

    // Create a new HTTP POST request
    req, err := http.NewRequest("POST", eventWebhookURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }
    req.Header.Set("Content-Type", "application/json")

    // Send the request
    resp, err := WebhookClient.Do(req)
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

	fmt.Println("Response From Internal Webhook:", resp.Status)
}

func sendBossKilledWebhook(bossName string) {

}
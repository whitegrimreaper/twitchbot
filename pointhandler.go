package main

import (
	"fmt"
)

// This class holds functions that hold the 'business logic' interfacing between the twitch commands
// and the actual db functions. This keeps some nice separation and compartmentalization
// It's also what we did at my only professional Golang gig so far so monkey see monkey implement

// Tries to add a set number of bosses to the queue
func executeBossAddition(idInt int, bossName string, numKills int)(response string) {
	respCode, respMessage, exists := doesUserExist(idInt)
	if respCode != 0 || respMessage != "" {
		response = ""
		return
	}
	if(!exists) {
		// User doesn't exist in DB so they probably haven't gained any points before
		// (we keep 0 point entries)
		response = "Error adding boss kills to log: you don't have any points yet!"
		return
	} else {
		// User exists, grab they points
		respCode, respMessage, points := findUserPoints(idInt)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = "Serious internal error: Id exists but points req failed!"
			return
		}
		respCode, respMessage, trueName := getBossTrueName(bossName)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = "Make sure to give a known name for the boss!"
			return
		}
		// Now also grab the boss info from the db
		respCode, respMessage, bossInfo := getBossWithName(trueName)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = "Error: boss name not known!"
			return
		}
		valid, reason := checkRequestIsValid(bossInfo, points, numKills)
		if !valid {
			response = reason
		}
	}
	return
}
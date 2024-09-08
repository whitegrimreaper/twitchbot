package main

import (
	"fmt"
	"strconv"
)

// This class holds functions that hold the 'business logic' interfacing between the twitch commands
// and the actual db functions. This keeps some nice separation and compartmentalization
// It's also what we did at my only professional Golang gig so far so monkey see monkey implement

// Tries to add a set number of bosses to the queue
func executeBossAddition(idInt int, bossName string, numKills int)(response string) {
	respCode, respMessage, exists := doesUserExist(idInt)
	if respCode != 0 || respMessage != "" {
		// TODO - Handle this correctly, the user should know if they try to add kills and have no 
		// points
		response = "You need to earn points before you add kills!"
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
			return
		}
		if bossInfo.BossRequestsFrozen != 0 {
			response = "Requests for that boss are frozen!"
			return
		}
		pointCost := (numKills * bossInfo.BossCost)
		fmt.Printf("numKills %d bosscost %d should cost %d\n", numKills, bossInfo.BossCost, pointCost)
		// Do the thing
		respCode, respMessage = writePointSpendEvent(idInt, pointCost)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = respMessage
			return
		}
		respCode, respMessage = addBossKillsQueue(idInt, bossInfo.BossID, numKills)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = respMessage
			return
		}
		respCode, respMessage = addBossKillsMain(idInt, bossInfo.BossID, numKills)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("%s\n", respMessage)
			response = respMessage
			return
		}
		response = strconv.Itoa(numKills) + " kills added for " + bossInfo.BossName + "!"
	}
	return
}

func executeBossRemoval(bossName string, numKills  int)(response string) {
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
	respCode, respMessage, queueId := findQueueEntryToRemoveFrom(bossInfo.BossID)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = "Error: no queue entries for this boss!"
		return
	}
	

	// do the thing
	respCode, respMessage = removeBossKillsQueue(queueId, numKills)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = "Error removing kills from the queue"
		return
	}

	// I'm just using add with a negative number here tbh
	respCode, respMessage = addBossKillsMain(queueId, bossInfo.BossID, numKills * -1)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = respMessage
		return
	}

	// check if we should delete the request entry
	respCode, respMessage = deleteQueueIfEmpty(queueId)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = respMessage
		return
	}
	response = "Success!"
	return
}

func handleCheckKills(userId int)(response string) {
	respCode, respMessage, entries := findUserQueueEntries(userId)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = respMessage
		return
	}
	str, err := getListOfBosses(entries);
	if err != nil {
		response = "Oops I errored again"
	}
	response = "You currently have " + str + " left!"
	return
}

func getBossCost(bossName string)(response string) {
	respCode, respMessage, trueName := getBossTrueName(bossName)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = "Make sure to give a known name for the boss!"
		return
	}

	respCode, respMessage, bossInfo := getBossWithName(trueName)
	if respCode != 0 || respMessage != "" {
		fmt.Printf("%s\n", respMessage)
		response = "Error: boss name not known!"
		return
	}
	response = bossName + " costs " + strconv.Itoa(bossInfo.BossCost) + " per kill"
	return
}
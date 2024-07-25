package main

import (
	"fmt"
	"strconv"
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func handleCommand(args []string, message twitch.PrivateMessage) (string, error) {
	var ret string
	commandName := strings.TrimPrefix(args[0], "!")
	switch {
	case commandName == command_dummyCommand:
		ret = "fricc you pigeon NoPigeons FinestPigeon"
	case commandName == command_checkPoints:
		userID, err := strconv.Atoi(message.User.ID)
		if err != nil {
			fmt.Printf("")
		}
		respCode, respMessage, points := findUserPoints(userID)
		if respCode != 0 || respMessage != "" {
			fmt.Printf("")
		}
		ret = "You've got " + strconv.Itoa(points) + " points to spend!"
	case commandName == command_addKills:
		if len(args) != 3 {
			ret = "Incorrect number of args! Usage: !addKills <boss name> <# of kills>"
			return ret, nil
		}
		userID, err := strconv.Atoi(message.User.ID)
		if err != nil {
			fmt.Printf("Strconv error handling addKills - error converting UserID %s\n", message.User.ID)
			ret = "Error converting UserID somehow"
			return ret, nil
		}
		bossName := args[1]
		numKills, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Strconv error handling addKills - error converting kill count %s\n", message.User.ID)
			ret = "Kill count needs to be a number"
			return ret, nil
		}
		// allow user to spend points to add kills to the queue
		ret = executeBossAddition(userID,bossName, numKills)
	case commandName == command_checkKills:
		// allow user to check how many kills they have left in the current queue
		ret = ""
	case commandName == command_help:
		ret = "Help is on the way! (help command is currently under construction)"
		// call separate "handle_help" function or smth
	}
	return ret, nil
}

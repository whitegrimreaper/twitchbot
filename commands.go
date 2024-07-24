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
		// allow user to spend points to add kills to the queue
		ret = ""
	case commandName == command_checkKills:
		// allow user to check how many kills they have left in the current queue
		ret = ""
	case commandName == command_help:
		ret = "Help is on the way! (help command is currently under construction)"
		// call separate "handle_help" function or smth
	}
	return ret, nil
}

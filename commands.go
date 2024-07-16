package main

import (
	"strings"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func handleCommand(args []string, message twitch.PrivateMessage) (string, error) {
	var ret string
	commandName := strings.TrimPrefix(args[0], "!")
	switch {
	case commandName == dummyCommand:
		ret = "fricc you pigeon NoPigeons"
	case commandName == checkPoints:
		// check user's current points
		ret = ""
	case commandName == addKills:
		// allow user to spend points to add kills to the queue
		ret = ""
	case commandName == checkKills:
		// allow user to check how many kills they have left in the current queue
		ret = ""
	}
	return ret, nil
}

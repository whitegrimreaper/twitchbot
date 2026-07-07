package main

import "strings"

// Generic stuff
const helpGeneral = "Commands: !check, !addKills, !checkKills, !checkCost - type !help <command> for more details"

// Per-command help strings, keyed by the same lowercased name used in consts.go
var helpText = map[string]string{
	command_checkPoints:   "!check - shows how many points you currently have",
	command_addKills:      "!addKills <boss name> <# of kills> - spend points to add kills to the queue. Example: !addKills Kril 5",
	command_checkKills:    "!checkKills - shows all your kills in the queue",
	command_checkCost:     "!checkCost <boss name> - shows how many points a single kill costs for that boss. Example: !checkCost Kril",
	command_howEarn:       "!howToEarn - explains how to earn boss points",
	command_removeKills:   "!removeKills - don't worry about this one",
	command_help:          "!help - now actually helps!. usage: !help <command>",
	command_dummyCommand:  "!pigeon - Insults pigeonmob",
	command_sandpitTurtle: "!sandpitturtle - will we ever see the end?",
}

func plsHelp(args []string) string {
	// with no args
	if len(args) < 2 {
		return helpGeneral
	}

	// with args
	commandName := strings.ToLower(strings.TrimPrefix(args[1], "!"))

	if text, ok := helpText[commandName]; ok {
		return text
	}

	return "Unknown command '" + args[1] + "'. "
}

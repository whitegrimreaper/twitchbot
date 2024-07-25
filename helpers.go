package main

import(
	"fmt"
	"os"
	helix "github.com/nicklaw5/helix/v2"
)

func printThingsFromEventSubResp(resp *helix.EventSubSubscriptionsResponse) {
	if(resp.ResponseCommon.StatusCode > 199 && resp.ResponseCommon.StatusCode  < 300) {
		// 200 class error means  we good
		fmt.Printf("Success!\n")
	} else if(resp.ResponseCommon.StatusCode > 300) {
		// counting any of this as error lol
		fmt.Printf("Failed!\n")
	} else {
		fmt.Printf("You have majorly fucked something up\n")
	}
}

func validateEnvVars(vars []string) {
	for idx, envVar := range vars {
		if envVar == "" {
			fmt.Printf("Var at idx %d is not set!\n", idx)
			os.Exit(0)
		}
	}
}

func bossNameToId(bossName string)(bossId int) {
	respCode, respMessage, boss := getBossWithName(bossName)
	if respCode != 0 || respMessage != "" {
		return 0
	}
	return boss.BossID
}

func checkRequestIsValid(bossInfo BossEntry, currPoints int, numKills int)(valid bool, reason string) {
	requestCost := numKills * bossInfo.BossCost;

	if requestCost > currPoints {
		return false, "User doesn't have enough points"
	} else {
		return true, ""
	}
}
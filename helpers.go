package main

import(
	"fmt"
	"os"
	"strings"
	helix "github.com/nicklaw5/helix/v2"
)

func printThingsFromEventSubResp(resp *helix.EventSubSubscriptionsResponse) {
	if(resp.ResponseCommon.StatusCode > 199 && resp.ResponseCommon.StatusCode  < 300) {
		// 200 class error means  we good
		fmt.Printf("Success!\n")
	} else if(resp.ResponseCommon.StatusCode > 300) {
		// counting any of this as error lol
		fmt.Printf("Failed!\n")
		fmt.Printf("%+v\n", resp)
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

// yoinked from fastbill/go-tiny-helpers because the import didn't work 
// (prob a golang version thing)
func ContainsStringCaseInsensitive(list []string, value string) bool {
	for _, item := range list {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}
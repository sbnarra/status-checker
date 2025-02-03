package main

import (
	"encoding/json"
	"fmt"
	"os"
	"status/internal/checker"
	"status/internal/model"
	"status/internal/slack"
	"status/internal/util"
)

func main() {
	var checkPath string
	if len(os.Args) > 1 {
		checkPath = os.Args[1]
	} else {
		checkPath = "checks.json"
	}

	slackHookUrl := os.Getenv("SLACK_HOOK_URL")

	if checker, err := checker.New(checkPath, onChecked(slackHookUrl)); err != nil {
		panic(err)
	} else {
		checker.Start()
		util.WaitUserTermination()
		checker.Stop()
	}
}

func onChecked(slackHookUrl string) checker.OnChecked {
	return func(name string, check model.Check, result model.CheckResult) {
		if result.CheckError != nil {
			if slackHookUrl != "" {
				slack.Send(slackHookUrl, name, check, result)
			} else {
				checkJson, _ := json.Marshal(check)
				resultJson, _ := json.Marshal(result)
				fmt.Printf("check=%s,result=%s\n", string(checkJson), string(resultJson))
			}
		}
	}
}

package main

import (
	"fmt"
	"os"
	"status/internal/checker"
	"status/internal/history"
	"status/internal/model"
	"status/internal/prometheus"
	"status/internal/server"
	"status/internal/slack"
	"status/internal/util"
)

func main() {
	var checkConfigPath = getCheckConfigPath()

	slackHookUrl := os.Getenv("SLACK_HOOK_URL")
	historyPath := os.Getenv("HISTORY_PATH")
	if historyPath == "" {
		historyPath = "history.yaml"
	}
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = ":8080"
	}
	prometheusEnabled := util.IsTrue(os.Getenv("PROMETHEUS_ENABLED"))

	onChecked := onChecked(historyPath, prometheusEnabled, slackHookUrl)
	if checker, err := checker.New(checkConfigPath, onChecked); err != nil {
		panic(err)
	} else {
		checker.Start()
		if err := server.Listen(bindAddr); err != nil {
			fmt.Println("listener error", err)
		}
		checker.Stop()
	}
}

func getCheckConfigPath() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	} else if checkPath := os.Getenv("CHECKS_PATH"); checkPath != "" {
		return checkPath
	} else {
		return "checks.yaml"
	}
}

func onChecked(historyPath string, prometheusEnabled bool, slackHookUrl string) checker.OnChecked {
	return func(name string, check model.Check, result model.CheckResult) {
		if err := history.Append(historyPath, name, result); err != nil {
			fmt.Println("failed to append history:", err)
		}

		if prometheusEnabled {
			prometheus.Publish(name, result)
		}

		if result.CheckError != nil {
			// only notify if webhook is configured and for recovered/failed checks
			if err := slack.Notify(slackHookUrl, name, check, result); err != nil {
				fmt.Println("failed to send slack message:", err)
			}
		}
	}
}

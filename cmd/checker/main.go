package main

import (
	"fmt"
	"status-checker/internal/checker"
	"status-checker/internal/config"
	"status-checker/internal/history"
	"status-checker/internal/prometheus"
	"status-checker/internal/server"
	"status-checker/internal/slack"
)

func main() {
	if checker, err := checker.New(onChecked); err != nil {
		panic(err)
	} else {
		checker.Start()
		defer checker.Stop()
	}

	if err := server.Listen(config.BindAddr); err != nil {
		fmt.Println("listener error", err)
	}
}

func onChecked(name string, check checker.Check, result checker.CheckResult) {
	if err := history.Append(name, result); err != nil {
		fmt.Println("failed to append history:", err)
	}

	prometheus.Publish(name, result)

	if result.CheckError != nil {
		// only notify if webhook is configured and for recovered/failed checks
		if err := slack.Notify(name, check, result); err != nil {
			fmt.Println("failed to send slack message:", err)
		}
	}
}

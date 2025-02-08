package main

import (
	"fmt"
	"os"
	"os/signal"
	"status-checker/internal/checker"
	"status-checker/internal/config"
	"status-checker/internal/history"
	"status-checker/internal/prometheus"
	"status-checker/internal/server"
	"status-checker/internal/slack"
	"syscall"
)

func main() {
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("config init error: %w", err))
	} else if checker, err := checker.New(onChecked); err != nil {
		panic(err)
	} else {
		checker.Start()
		defer checker.Stop()

		if config.ServerEnabled {
			if err := server.Listen(config.BindAddr); err != nil {
				fmt.Println("listener error", err)
			}
		} else {
			exitSignal := make(chan os.Signal, 1)
			signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
			<-exitSignal
		}
	}
}

func onChecked(name string, check checker.Check, result checker.Result) {
	if err := history.Append(name, result); err != nil {
		fmt.Println("failed to append history:", err)
	}
	if result.CheckError != nil {
		if err := slack.Notify(name, check, result); err != nil {
			fmt.Println("failed to send slack message:", err)
		}
	}
	prometheus.Publish(name, result)
}

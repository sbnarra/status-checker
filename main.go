package main

import (
	"fmt"
	"os"
	"os/signal"
	"status-checker/internal/checker"
	"status-checker/internal/config"
	"status-checker/internal/history"
	"status-checker/internal/log"
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
		log.Error("failed to append history: %s", err)
	}
	prometheus.Publish(name, result)
	if result.CheckError == nil {
		return
	}
	if err := slack.Notify(name, check, result); err != nil {
		log.Error("failed to send slack message: %s", err)
	}
}

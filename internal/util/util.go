package util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitUserTermination() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

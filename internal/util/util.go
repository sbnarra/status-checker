package util

import (
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
)

func WaitUserTermination() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func IsTrue(val string) bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(val))
}

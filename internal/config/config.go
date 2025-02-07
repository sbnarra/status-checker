package config

import (
	"os"
	"slices"
	"strings"
)

var BindAddr = fallback("BIND_ADDR", ":8000")
var CheckPath = args("CHECKS_PATH", "config/checks.yaml")
var HistoryPath = fallback("HISTORY_PATH", "history/history.yaml")
var SlackHookUrl = fallback("SLACK_HOOK_URL", "https://hooks.slack.com/services/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")
var PrometheusEnabled = enabled(os.Getenv("PROMETHEUS_ENABLED"))

func args(key string, val string) string {
	if len(os.Args) > 1 {
		return os.Args[1]
	} else {
		return fallback(key, val)
	}
}

func fallback(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func enabled(val string) bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(val))
}

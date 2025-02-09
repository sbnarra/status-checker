package config

import (
	"fmt"
	"os"
)

var BindAddr string
var CheckPaths []string
var HistoryDir string
var MinHistory int
var HistoryCheckSizeLimit uintptr
var SlackHookUrl string
var ServerEnabled bool
var PrometheusEnabled bool
var Debug bool

func Init() error {
	var err error
	// checker config
	if len(os.Args) > 1 {
		CheckPaths = os.Args[1:]
	} else if CheckPaths, err = withStrArrDefault("CHECKS_PATH", []string{"./config/checks.yaml"}); err != nil {
		return err
	}
	// server
	if ServerEnabled, err = withBoolDefault("SERVER_ENABLED", true); err != nil {
		return err
	} else if BindAddr, err = withStrDefault("BIND_ADDR", ":8000"); err != nil {
		return err
	} else if Debug, err = withBoolDefault("DEBUG", true); err != nil {
		return err
	}
	// history
	if HistoryDir, err = withStrDefault("HISTORY_DIR", "./history"); err != nil {
		return err
	} else if MinHistory, err = withIntDefault("MIN_HISTORY", 100); err != nil {
		return err
	} else if historyCheckSizeLimit, err := withStrDefault("HISTORY_CHECK_SIZE_LIMIT", "10MB"); err != nil {
		return err
	} else if HistoryCheckSizeLimit, err = toBytes(historyCheckSizeLimit); err != nil {
		return fmt.Errorf("HISTORY_CHECK_SIZE_LIMIT byte conversion error: %w", err)
	}
	// monitoring/notifications
	if PrometheusEnabled, err = withBoolDefault("PROMETHEUS_ENABLED", false); err != nil {
		return err
	} else if SlackHookUrl, err = withStrDefault("SLACK_HOOK_URL", "https://hooks.slack.com/services/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"); err != nil {
		return err
	}
	return nil
}

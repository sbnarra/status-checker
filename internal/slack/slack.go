package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"status-checker/internal/checker"
	"status-checker/internal/config"
)

func Notify(name string, check checker.Check, result *checker.Result) error {
	if config.SlackHookUrl == "" {
		return fmt.Errorf("missing slack hook url")
	}
	jsonData, err := json.Marshal(map[string]string{"text": Message(name, check, result)})
	if err != nil {
		return fmt.Errorf("failed to encode slack payload: %w", err)
	}

	_, err = http.Post(config.SlackHookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post to slack: %w", err)
	}
	return nil
}

func Message(name string, check checker.Check, result *checker.Result) string {
	return commandMarkdown("Check", name, check, result.Check) +
		commandMarkdown("Recovery", name, check, result.Recover) +
		commandMarkdown("Re-Check", name, check, result.ReCheck)
}

func commandMarkdown(stage string, name string, check checker.Check, cmd *checker.CmdResult) string {
	if cmd == nil {
		return fmt.Sprintf("\n*%s Skipped: _%s_*\n_Command:_ `%s`", stage, name, check.Command)
	} else if cmd.Error != nil {
		messagePrefix := "\n*%s Error: _%s_*\n_Command:_ `%s`\n_Error:_ `%s`"
		if cmd.Output == "" {
			return fmt.Sprintf(messagePrefix, stage, name, cmd.Command, cmd.Error)
		}
		return fmt.Sprintf(messagePrefix+"\n```\n%s\n```", stage, name, cmd.Command, cmd.Error, cmd.Output)
	}
	return fmt.Sprintf("\n*%s Success: _%s_*\n_Command:_ `%s`\n```\n%s\n```", stage, name, cmd.Command, cmd.Output)
}

package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"status/internal/model"
)

func Send(slackHookUrl string, name string, check model.Check, result model.CheckResult) error {
	jsonData, err := json.Marshal(map[string]string{"text": Message(name, check, result)})
	if err != nil {
		return fmt.Errorf("failed to encode slack payload: %w", err)
	}

	_, err = http.Post(slackHookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to post to slack: %w", err)
	}
	return nil
}

func Message(name string, check model.Check, result model.CheckResult) string {
	return commandMarkdown("Check", name, &check.Command, &result.CheckOutput, result.CheckError) +
		commandMarkdown("Recovery", name, check.Recover, result.RecoverOutput, result.RecoverError) +
		commandMarkdown("Re-Check", name, &check.Command, result.RecheckOutput, result.RecheckError)
}

func commandMarkdown(stage string, name string, command *string, output *string, err error) string {
	if command == nil {
		noCommand := "No Command"
		command = &noCommand
	}

	if err != nil {
		messagePrefix := "\n*%s Error: _%s_*\n_Command:_ `%s`\n_Error:_ `%s`"
		if output == nil {
			return fmt.Sprintf(messagePrefix, stage, name, *command, err)
		}
		return fmt.Sprintf(messagePrefix+"\n```\n%s\n```", stage, name, *command, err, *output)
	} else if output == nil {
		return fmt.Sprintf("\n*%s Skipped: _%s_*\n_Command:_ `%s`", stage, name, *command)

	}
	return fmt.Sprintf("\n*%s Success: _%s_*\n_Command:_ `%s`\n```\n%s\n```", stage, name, *command, *output)
}

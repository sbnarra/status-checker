package checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"status/internal/model"
	"strings"

	"github.com/robfig/cron/v3"
)

type OnChecked func(string, model.Check, model.CheckResult)

func New(checkPath string, onChecked OnChecked) (*cron.Cron, error) {
	checks, err := readConfig(checkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	checker := startChecking(checks, onChecked)
	return checker, nil
}

func startChecking(checks map[string]model.Check, onChecked OnChecked) *cron.Cron {
	c := cron.New()
	for name, check := range checks {
		registerCheck(name, check, onChecked, c)
	}
	return c
}

func registerCheck(name string, check model.Check, onChecked OnChecked, c *cron.Cron) {
	var schedule string
	if check.Schedule == nil {
		schedule = "* * * * *"
	} else {
		schedule = *check.Schedule
	}

	fmt.Printf("check '%s' has schedule '%s' using '%s'\n", name, schedule, check.Command)
	result := runCheck(check)
	onChecked(name, check, result)
	c.AddFunc(schedule, func() {
		result := runCheck(check)
		onChecked(name, check, result)
	})
}

func runCheck(check model.Check) model.CheckResult {
	checkStdout, checkError := runCmd(check.Command)
	result := model.CheckResult{
		CheckOutput: checkStdout,
		CheckError:  checkError,
	}

	if result.CheckError == nil {
		return result
	} else if check.Recover == nil {
		result.RecoverError = errors.New("no recovery command")
		return result
	}

	recoverStdout, recoverErr := runCmd(*check.Recover)
	result.RecoverOutput = &recoverStdout
	result.RecoverError = recoverErr
	if result.RecoverError != nil {
		return result
	}

	recheckStdout, recheckErr := runCmd(check.Command)
	result.RecheckOutput = &recheckStdout
	result.RecheckError = recheckErr
	return result
}

func runCmd(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func readConfig(checkPath string) (map[string]model.Check, error) {
	var checks map[string]model.Check
	if checksFile, err := os.Open(checkPath); err != nil {
		return nil, err
	} else if checksData, err := io.ReadAll(checksFile); err != nil {
		return nil, err
	} else if err := json.Unmarshal(checksData, &checks); err != nil {
		return checks, err
	} else {
		return checks, nil
	}
}

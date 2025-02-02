package checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"status/internal/model"

	"github.com/robfig/cron/v3"
)

type onChecked func(string, model.Check, model.CheckResult)

func New(checkPath string, onChecked onChecked) (*cron.Cron, error) {
	checks, err := readConfig(checkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	checker := startChecking(checks, onChecked)
	return checker, nil
}

func startChecking(checks map[string]model.Check, onChecked onChecked) *cron.Cron {
	c := cron.New()
	for name, check := range checks {
		registerCheck(name, check, onChecked, c)
	}
	return c
}

func registerCheck(name string, check model.Check, onChecked onChecked, c *cron.Cron) {
	var schedule string
	if check.Schedule == nil {
		schedule = "* * * * *"
	} else {
		schedule = *check.Schedule
	}

	fmt.Printf("check '%s' has schedule '%s' using '%s'\n", name, schedule, check.Command)
	c.AddFunc(schedule, func() {
		result := runCheck(check)
		onChecked(name, check, result)
	})
}

func runCheck(check model.Check) model.CheckResult {
	out, err := runCmd(check.Command)
	result := model.CheckResult{
		CheckOutput: out,
		CheckError:  &err,
	}

	if result.CheckError == nil {
		return result
	} else if check.Recover == nil {
		err = errors.New("no recovery command")
		result.RecoverError = &err
		return result
	}

	out, err = runCmd(*check.Recover)
	result.RecoverOutput = &out
	result.RecoverError = &err
	if result.RecoverError != nil {
		return result
	}

	out, err = runCmd(check.Command)
	result.RecheckOutput = &out
	result.RecheckError = &err
	return result
}

func runCmd(command string) (string, error) {
	fmt.Println("executing:", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return string(out), err
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

package checker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"status/internal/model"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
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
		fmt.Printf("check '%s' has schedule '%s' using '%s'\n", name, check.Schedule, check.Command)
		CheckNow(name, check, onChecked)
		c.AddFunc(check.Schedule, func() {
			CheckNow(name, check, onChecked)
		})
	}
	return c
}

func CheckNow(name string, check model.Check, onChecked OnChecked) {
	result := runCheck(check)
	result.Completed = time.Now()
	onChecked(name, check, result)
}

func runCheck(check model.Check) model.CheckResult {
	result := model.CheckResult{
		Started: time.Now(),
	}

	checkStdout, checkError := runCmd(check.Command)
	result.CheckOutput = checkStdout
	result.CheckError = errToStr(checkError)

	if result.CheckError == nil {
		result.Status = model.CheckSuccess
		return result
	} else if check.Recover == nil {
		missingRecovery := "no recovery command"
		result.RecoverError = &missingRecovery
		result.Status = model.CheckFailed
		return result
	}

	recoverStdout, recoverErr := runCmd(*check.Recover)
	result.RecoverOutput = &recoverStdout
	result.RecoverError = errToStr(recoverErr)
	if result.RecoverError != nil {
		result.Status = model.CheckFailed
		return result
	}

	recheckStdout, recheckErr := runCmd(check.Command)
	result.RecheckOutput = &recheckStdout
	result.RecheckError = errToStr(recheckErr)
	result.Status = model.CheckRecovered

	return result
}

func errToStr(err error) *string {
	if err == nil {
		return nil
	}
	errStr := fmt.Sprintf("%s", err)
	return &errStr
}

func runCmd(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func readConfig(checkPath string) (map[string]model.Check, error) {
	checks := make(map[string]model.Check)

	checksFile, err := os.Open(checkPath)
	if err != nil {
		return nil, err
	}
	defer checksFile.Close()

	if checksData, err := io.ReadAll(checksFile); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(checksData, &checks); err != nil {
		return checks, err
	}
	return checks, nil
}

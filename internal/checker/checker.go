package checker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"status-checker/internal/config"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

type OnChecked func(string, Check, Result)

func New(onChecked OnChecked) (*cron.Cron, error) {
	checks, err := Config()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	checker := startChecking(checks, onChecked)
	return checker, nil
}

func startChecking(checks map[string]Check, onChecked OnChecked) *cron.Cron {
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	for name, check := range checks {

		fmt.Printf("check '%s' has schedule '%s' using '%s'\n", name, check.Schedule, check.Command)
		result := runCheck(check)
		onChecked(name, check, result)

		c.AddFunc(check.Schedule, func() {
			fmt.Println("running check", name)
			result := runCheck(check)
			onChecked(name, check, result)
		})
	}
	return c
}

func runCheck(check Check) Result {
	result := Result{
		Command: check.Command,
		Recover: check.Recover,
		Started: time.Now(),
		Status:  "Failed",
	}
	defer func() { result.Completed = time.Now() }()

	checkStdout, checkError := runCmd(check.Command)
	result.CheckOutput = checkStdout
	result.CheckError = errToStr(checkError)
	if result.CheckError == nil {
		result.Status = "Success"
		return result
	} else if check.Recover == nil {
		missingRecovery := "no recovery command"
		result.RecoverError = &missingRecovery
		return result
	}

	recoverStdout, recoverErr := runCmd(*check.Recover)
	result.RecoverOutput = &recoverStdout
	result.RecoverError = errToStr(recoverErr)
	if result.RecoverError != nil {
		return result
	}

	recheckStdout, recheckErr := runCmd(check.Command)
	result.RecheckOutput = &recheckStdout
	result.RecheckError = errToStr(recheckErr)
	if result.RecheckError != nil {
		return result
	}
	result.Status = "Recovered"
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

var checks map[string]Check

func Config() (map[string]Check, error) {
	if checks != nil {
		return checks, nil
	}

	loadedChecks := map[string]Check{}
	for _, path := range config.CheckPaths {
		if loaded, err := load(path); err != nil {
			return loadedChecks, err
		} else {
			for name, check := range loaded {
				loadedChecks[name] = check
			}
		}
	}
	checks = loadedChecks
	return checks, nil
}

func load(path string) (map[string]Check, error) {
	checksFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer checksFile.Close()

	loaded := make(map[string]Check)
	if checksData, err := io.ReadAll(checksFile); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(checksData, &loaded); err != nil {
		return loaded, err
	}
	return loaded, nil
}

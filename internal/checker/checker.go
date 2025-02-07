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

type OnChecked func(string, Check, CheckResult)

func New(onChecked OnChecked) (*cron.Cron, error) {
	checks, err := ReadConfig()
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
		CheckNow(name, check, onChecked)
		c.AddFunc(check.Schedule, func() {
			fmt.Println("running check", name)
			CheckNow(name, check, onChecked)
		})
	}
	return c
}

func CheckNow(name string, check Check, onChecked OnChecked) {
	result := runCheck(check)
	result.Completed = time.Now()
	onChecked(name, check, result)
}

func runCheck(check Check) CheckResult {
	result := CheckResult{
		Started: time.Now(),
	}

	checkStdout, checkError := runCmd(check.Command)
	result.CheckOutput = checkStdout
	result.CheckError = errToStr(checkError)

	if result.CheckError == nil {
		result.Status = "Success"
		return result
	}

	result.Status = "Failed"
	if check.Recover == nil {
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

func ReadConfig() (map[string]Check, error) {
	checks := make(map[string]Check)

	checksFile, err := os.Open(config.CheckPath)
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

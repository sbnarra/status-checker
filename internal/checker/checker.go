package checker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"status-checker/internal/config"
	"status-checker/internal/log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

type CheckCallback func(string, Check, *Result)

func New(callback CheckCallback) (*cron.Cron, error) {
	if checks, err := Config(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	} else if checker, err := runChecks(checks, callback); err != nil {
		return nil, fmt.Errorf("failed to start checker: %w", err)
	} else {
		return checker, nil
	}
}

func runChecks(checks map[string]Check, callback CheckCallback) (*cron.Cron, error) {
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	for name, check := range checks {

		log.Info("check '%s' has schedule '%s' using '%s'", name, check.Schedule, check.Command)
		runCheckNow := func() {
			if err := runCheck(name, check, callback); err != nil {
				log.Error("failed to run check '%s': %s", name, err)
			}
		}
		go runCheckNow()
		if _, err := c.AddFunc(check.Schedule, runCheckNow); err != nil {
			return nil, err
		}
	}
	return c, nil
}

var checkRunLocks sync.Map

func runCheck(name string, check Check, callback CheckCallback) error {
	lock, _ := checkRunLocks.LoadOrStore(name, &sync.Mutex{})
	lock.(*sync.Mutex).Lock()
	defer lock.(*sync.Mutex).Unlock()

	log.Debug("running check '%s'", name)
	result := &Result{
		Status:  StatusRunning,
		Started: time.Now(),
	}
	callback(name, check, result)
	defer callback(name, check, result)

	finalise := func(status Status) {
		result.Completed = time.Now()
		result.Status = status
	}
	onError := func(status Status, stage string, checkErr error, err error) error {
		finalise(status)
		return fmt.Errorf("check=%w: %s=%w", checkErr, stage, err)
	}

	result.Check = &CmdResult{}
	if checkErr := runCmd(&check.Command, result.Check); checkErr != nil {
		result.Recover = &CmdResult{}
		if err := runCmd(check.Recover, result.Recover); err != nil {
			return onError(StatusFailed, "recover", checkErr, err)
		}
		result.ReCheck = &CmdResult{}
		if err := runCmd(&check.Command, result.ReCheck); err != nil {
			return onError(StatusFailed, "recheck", checkErr, err)
		}
	} else {
		finalise(StatusSuccess)
	}
	finalise(StatusRecovered)
	return nil
}

func runCmd(command *string, result *CmdResult) error {
	result.Status = StatusRunning
	result.Started = time.Now()
	finalise := func(status Status, err error) error {
		result.Completed = time.Now()
		result.Status = status
		result.Error = errToStr(err)
		return err
	}

	if command == nil {
		return finalise(StatusFailed, fmt.Errorf("missing command"))
	}
	result.Command = *command

	cmd := exec.Command("bash", "-c", result.Command)
	cmd.Env = os.Environ()

	if stdout, err := cmd.StdoutPipe(); err != nil {
		return finalise(StatusFailed, err)
	} else if stderr, err := cmd.StderrPipe(); err != nil {
		return finalise(StatusFailed, err)
	} else {
		capture := func(reader io.ReadCloser) {
			scanner := bufio.NewScanner(reader)
			scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				result.Output += scanner.Text()
			}
		}
		go capture(stdout)
		go capture(stderr)
	}

	if err := cmd.Run(); err != nil {
		return finalise(StatusFailed, err)
	}
	return finalise(StatusSuccess, nil)
}

func errToStr(err error) *string {
	if err == nil {
		return nil
	}
	errStr := fmt.Sprintf("%s", err)
	return &errStr
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

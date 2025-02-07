package history

import (
	"io"
	"os"
	"path/filepath"
	"status-checker/internal/checker"
	"status-checker/internal/config"
	"sync"

	"gopkg.in/yaml.v3"
)

var mu sync.Mutex

func Append(name string, result checker.CheckResult) error {
	mu.Lock()
	defer mu.Unlock()

	history, err := Read()
	if err != nil {
		return err
	} else if _, ok := history[name]; !ok {
		history[name] = []checker.CheckResult{}
	}
	history[name] = append(history[name], result)

	if historyContent, err := yaml.Marshal(history); err != nil {
		return err
	} else if err := os.MkdirAll(filepath.Dir(config.HistoryPath), os.ModePerm); err != nil {
		return err
	} else if err := os.WriteFile(config.HistoryPath, historyContent, 0644); err != nil {
		return err
	}
	return nil
}

func Read() (map[string][]checker.CheckResult, error) {
	history := make(map[string][]checker.CheckResult)

	historyFile, err := os.Open(config.HistoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return history, nil
		}
		return nil, err
	}
	defer historyFile.Close()

	if historyData, err := io.ReadAll(historyFile); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(historyData, &history); err != nil {
		return history, err
	}
	return history, nil
}

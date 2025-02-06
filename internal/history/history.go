package history

import (
	"io"
	"os"
	"status/internal/model"
	"sync"

	"gopkg.in/yaml.v3"
)

var mu sync.Mutex

func Append(historyPath string, name string, result model.CheckResult) error {
	mu.Lock()
	defer mu.Unlock()

	history, err := read(historyPath)
	if err != nil {
		return err
	} else if _, ok := history[name]; !ok {
		history[name] = []model.CheckResult{}
	}
	history[name] = append(history[name], result)

	if historyContent, err := yaml.Marshal(history); err != nil {
		return err
	} else if err := os.WriteFile(historyPath, historyContent, 0644); err != nil {
		return err
	}
	return nil
}

func read(historyPath string) (map[string][]model.CheckResult, error) {
	history := make(map[string][]model.CheckResult)

	historyFile, err := os.Open(historyPath)
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

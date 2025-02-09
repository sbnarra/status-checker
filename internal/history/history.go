package history

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"status-checker/internal/checker"
	"status-checker/internal/config"
	"sync"
	"unsafe"
)

var mu sync.Mutex
var historys = map[string][]checker.Result{}

func Append(name string, result checker.Result) error {
	mu.Lock() // TODO: only lock per-name?
	defer mu.Unlock()

	history, err := Get(name)
	if err != nil {
		return err
	}
	history = append(history, result)

	for memoryUsage(history) > config.HistoryCheckSizeLimit {
		if length := len(history); length > config.MinHistory {
			history = history[:length-1]
		} else {
			break
		}
	}
	historys[name] = history

	if historyContent, err := json.Marshal(history); err != nil {
		return err
	} else if err := os.MkdirAll(config.HistoryDir, os.ModePerm); err != nil {
		return err
	} else if err := os.WriteFile(historyFile(name), historyContent, 0644); err != nil {
		return err
	}
	return nil
}

func memoryUsage(arr []checker.Result) uintptr {
	var usage uintptr = 0
	for _, item := range arr {
		usage += unsafe.Sizeof(item)
	}
	return usage
}

func Get(name string) ([]checker.Result, error) {
	if history, ok := historys[name]; ok {
		return history, nil
	}

	if history, err := load(name); err != nil {
		return nil, err
	} else {
		historys[name] = history
		return history, nil
	}
}

func load(name string) ([]checker.Result, error) {
	history := []checker.Result{}

	if historyFile, err := os.Open(historyFile(name)); err != nil {
		if os.IsNotExist(err) {
			return history, nil
		}
		return nil, err
	} else {
		defer historyFile.Close()

		if historyData, err := io.ReadAll(historyFile); err != nil {
			return nil, err
		} else {
			err := json.Unmarshal(historyData, &history)
			return history, err
		}
	}
}

func historyFile(name string) string {
	return filepath.Join(config.HistoryDir, name+".json")
}

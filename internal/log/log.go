package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"status-checker/internal/config"
	"time"
)

func Info(msg string, vals ...any) {
	println("INFO", msg, vals...)
}

func Debug(msg string, vals ...any) {
	if config.Debug {
		println("DEBUG", msg, vals...)
	}
}

func Error(msg string, vals ...any) {
	println("ERROR", msg, vals...)
}

func println(level string, msg string, vals ...any) {
	level = "[" + level + "]"
	timestamp := time.Now().Format("2006/01/02 15:04:05.999")
	msg = fmt.Sprintf(msg, vals...)
	fmt.Printf("%s\t%s\t%s: %s\n", timestamp, level, caller(), msg)
}

func caller() string {
	if pc, _, line, ok := runtime.Caller(3); ok {
		funcForPc := runtime.FuncForPC(pc)
		file := filepath.Base(funcForPc.Name())
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}

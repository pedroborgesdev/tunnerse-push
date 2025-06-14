package logger

import (
	"fmt"
	"runtime"
	"time"

	"tunnerse/debug"
)

type LogDetail struct {
	Key   string
	Value interface{}
}

func Log(level string, message string, details []LogDetail) {
	if !debug.DebugConfig.Debug {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	color := getLevelColor(level)
	reset := "\033[0m"

	fmt.Printf("%s %s%s:%d%s \n↳ %s%s%s - %s\n",
		timestamp, color, file, line, reset,
		color, level, reset, message)

	for _, detail := range details {
		fmt.Printf("	%s%s%s: %v\n", color, detail.Key, reset, detail.Value)
	}
}

func getLevelColor(level string) string {
	switch level {
	case "DEBUG":
		return "\033[36m"
	case "INFO":
		return "\033[32m"
	case "WARN":
		return "\033[33m"
	case "ERROR":
		return "\033[31m"
	default:
		return "\033[35m"
	}
}

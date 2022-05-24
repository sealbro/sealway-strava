package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func Tracef(format string, v ...interface{}) {
	if os.Getenv("LOG_LEVEL") == "TRACE" {
		Printf("TRACE", format, v...)
	}
}
func Trace(message string) {
	if os.Getenv("LOG_LEVEL") == "TRACE" {
		Printf("TRACE", message)
	}
}

func InfoJson(jsonObj interface{}) {
	jsonBin, err := json.Marshal(jsonObj)
	if err != nil {
		fmt.Println(err)
		return
	}
	Printf("INFO", string(jsonBin))
}

func Infof(format string, v ...interface{}) {
	Printf("INFO", format, v...)
}
func Info(message string) {
	Printf("INFO", message)
}

func Warnf(format string, v ...interface{}) {
	Printf("WARN", format, v...)
}
func Warn(message string) {
	Printf("WARN", message)
}

func Errorf(format string, v ...interface{}) {
	Printf("ERROR", format, v...)
}
func Error(message string) {
	Printf("ERROR", message)
}

func Fatalf(format string, v ...interface{}) {
	Printf("FATAL", format, v...)
}
func Fatal(message string) {
	Printf("FATAL", message)
}

type LogEvent struct {
	Timestamp string `json:"@timestamp"`
	Level     string `json:"level"`
	Caller    string `json:"caller"`
	Message   string `json:"message"`
}

func Printf(level string, format string, v ...interface{}) {
	event := &LogEvent{
		Caller:    getCaller(2),
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   fmt.Sprintf(format, v...),
	}

	jsonText, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	dest := os.Stdout

	if level == "ERROR" || level == "FATAL" {
		dest = os.Stderr
	}

	_, err = dest.WriteString(fmt.Sprintf("%s\n", jsonText))
	if err != nil {
		panic(err)
	}

	if level == "FATAL" {
		os.Exit(1)
	}
}

func getCaller(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth + 1)
	if !ok {
		file = "???"
		line = 0
	}

	split := strings.Split(file, "/")
	file = split[len(split)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

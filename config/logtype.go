package config

import (
	"log/slog"
	"strings"
)

type LogType string

const (
	LogTypeJSON LogType = "json"
	LogTypeText LogType = "text"
)

func ParseLogType(v string) (LogType, error) {
	switch strings.ToLower(v) {
	case "json":
		return LogTypeJSON, nil
	case "text":
		return LogTypeText, nil
	default:
		slog.Warn("Invalid log type, use default type: json")
		return LogTypeJSON, nil
	}
}

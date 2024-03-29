package logger

import (
	"github.com/hckops/hckctl/pkg/util"
)

type LogLevel uint8

const (
	DebugLogLevel LogLevel = iota
	InfoLogLevel
	WarningLogLevel
	ErrorLogLevel
)

var levels = map[LogLevel]string{
	DebugLogLevel:   "debug",
	InfoLogLevel:    "info",
	WarningLogLevel: "warning",
	ErrorLogLevel:   "error",
}

func (l LogLevel) String() string {
	return levels[l]
}

func LevelValues() []string {
	return util.IotaToValues[LogLevel](levels)
}

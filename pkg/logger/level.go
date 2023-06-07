package logger

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
	var values []string
	for _, level := range levels {
		values = append(values, level)
	}
	return values
}

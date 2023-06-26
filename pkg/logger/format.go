package logger

type LogFormat uint8

const (
	JsonLogFormat LogFormat = iota
	TextLogFormat
)

var formats = map[LogFormat]string{
	JsonLogFormat: "json",
	TextLogFormat: "text",
}

func (l LogFormat) String() string {
	return formats[l]
}

func FormatValues() []string {
	return toValues[LogFormat](formats)
}

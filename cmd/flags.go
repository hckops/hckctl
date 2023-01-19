package cmd

import "github.com/spf13/viper"

const (
	LogLevelFlag = "log-level"
)

type Flags struct {
	LogLevel string
}

func NewFlags() *Flags {
	return &Flags{
		LogLevel: viper.GetString(LogLevelFlag),
	}
}

package common

const NoneFlagShortHand = ""

type GlobalCmdOptions struct {
	LogLevel       string // TODO enum
	InternalConfig *Config
}

type CommonCmdOptions struct {
	ConfigRef *Config
}

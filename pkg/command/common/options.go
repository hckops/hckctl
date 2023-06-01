package common

const NoneFlagShortHand = ""

type GlobalCmdOptions struct {
	LogLevel       string // TODO enum
	InternalConfig *ConfigV1
}

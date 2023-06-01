package common

// TODO delete
type GlobalCmdOptions struct {
	LogLevel       string // TODO enum
	InternalConfig *Config
}

// TODO ConfigRef
type CommonCmdOptions struct {
	Config *Config
}

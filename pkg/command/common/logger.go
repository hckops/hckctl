package common

import (
	"fmt"
	"log"
)

const (
	LogLevelFlag = "log-level"
)

func InitFileLogger(global *GlobalCmdOptions) {
	log.Println(fmt.Sprintf("LEVEL=%s", global.logLevel))
}

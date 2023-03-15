package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultDirectoryMod os.FileMode = 0755
	DefaultFileMod      os.FileMode = 0600
)

const (
	CliName       string = "hckctl"
	ConfigDir     string = "hck"
	ConfigNameEnv string = "HCK_CONFIG" //  overrides .config/hck/config.yml
)

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.log", CliName, GetUserOrDie()))

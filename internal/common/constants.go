package common

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
	ProjectName   string = "hckops"
	CliName       string = "hckctl"
	ConfigDir     string = "hck"
	ConfigNameEnv string = "HCK_CONFIG" //  overrides .config/hck/config.yml
)

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.log", CliName, GetUserOrDie()))

// TODO move common+model > config
// TODO add cli version: git + timestamp

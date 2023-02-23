package common

import "os"

const (
	ApiUrl            string = "https://api.hckops.com"
	SchemaUrl         string = "https://schema.hckops.com"
	MegalopolisRawUrl string = "https://raw.githubusercontent.com/hckops/megalopolis"
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
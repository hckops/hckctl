package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/hckops/hckctl/internal/common"
)

type CliConfig struct {
	log *LogConfig
	box *BoxConfig
}

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.log", common.CliName, common.GetUserOrDie()))

type LogConfig struct {
	LogLevel string
	LogFile  string
}

type BoxConfig struct {
	revision string
	kube     *KubeConfig
}

type KubeConfig struct {
	namespace     string
	k8sConfigPath string
}

func NewCliConfig() *CliConfig {
	return &CliConfig{
		log: &LogConfig{
			LogLevel: "info",
			LogFile:  DefaultLogFile,
		},
		box: &BoxConfig{
			revision: "main",
			kube: &KubeConfig{
				namespace:     "labs",
				k8sConfigPath: "~/.kube/config",
			},
		},
	}
}

// TODO
func (config *CliConfig) Setup() {
	home := common.ConfigHome()
	viper.AddConfigPath(home)
	viper.SetConfigName(fmt.Sprintf(".%", common.ConfigName))
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix(common.ConfigName)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("error reading config: %w", err))
	}
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/internal/common"
)

type CliConfig struct {
	Revision string    `yaml:"revision"`
	Box      BoxConfig `yaml:"box"`
	Log      LogConfig `yaml:"log"`
}

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.log", common.CliName, common.GetUserOrDie()))

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type BoxConfig struct {
	Kube KubeConfig `yaml:"kube"`
}

type KubeConfig struct {
	Namespace  string `yaml:"namespace"`
	ConfigPath string `yaml:"configPath"`
}

func newCliConfig() *CliConfig {
	return &CliConfig{
		Revision: "main",
		Box: BoxConfig{
			Kube: KubeConfig{
				Namespace:  "labs",
				ConfigPath: "~/.kube/config",
			},
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: DefaultLogFile,
		},
	}
}

func InitCliConfig() *CliConfig {

	cliConfig := newCliConfig()
	fmt.Printf("%+v\n", cliConfig)
	fmt.Printf("%#v\n", cliConfig)

	cliConfig.Print()

	return cliConfig
}

func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "prints current configurations",
		Run: func(cmd *cobra.Command, args []string) {

			// TODO
			fmt.Println(viper.AllSettings())
		},
	}
}

func (config *CliConfig) Print() {
	value, err := common.ToYaml(&config)
	if err != nil {
		log.Warn().Msg("invalid config")
	}
	fmt.Println(value)
}

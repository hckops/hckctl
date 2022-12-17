package cmd

import "github.com/spf13/viper"

const (
	ServerUrl = "server-url"
	Token     = "token"
	Local     = "local"
)

type globalFlags struct {
	serverUrl string
	token     string
	local     bool
}

func GlobalFlags() *globalFlags {
	return &globalFlags{
		serverUrl: viper.GetString(ServerUrl),
		token:     viper.GetString(Token),
		local:     viper.GetBool(Local),
	}
}

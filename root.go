package wgcf

import (
	"errors"
	"log"

	"github.com/KazumiLine/wgcf/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
}

var unsupportedConfigError viper.UnsupportedConfigError

func initConfig() {
	initConfigDefaults()
	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix("WGCF")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); errors.As(err, &unsupportedConfigError) {
		log.Fatal(err)
	} else {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initConfigDefaults() {
	viper.SetDefault(config.DeviceId, "")
	viper.SetDefault(config.AccessToken, "")
	viper.SetDefault(config.PrivateKey, "")
	viper.SetDefault(config.LicenseKey, "")
}

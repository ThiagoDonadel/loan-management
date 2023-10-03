package infra

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {

	//get profile for configuration
	viper.SetDefault("profile", "local")
	viper.BindEnv("profile", "MY_PROFILE")

}

func LoadConfigurationFromFile() error {

	configFile := fmt.Sprintf("config-%v", viper.GetString("profile"))

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil

}

package infra

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitEnv() {

	//bind env variables and default values for then
	viper.SetDefault("profile", "local")
	viper.BindEnv("profile", "LOAN_API_CONFIG_PROFILE")

	viper.SetDefault("config-profile", "../configs")
	viper.BindEnv("config-profile", "LOAN_API_CONFIG_FOLDER")

}

func LoadConfigurationFromFile() error {

	configFile := fmt.Sprintf("config-%v", viper.GetString("profile"))

	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("config-profile"))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

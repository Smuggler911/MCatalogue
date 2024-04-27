package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DB          string `mapstructure:"DB_STORAGE"`
	ExternalApi string `mapstructure:"EXTERNAL_API"`
}

func LoadConfig() (config Config, err error) {

	viper.SetConfigFile("././.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&config)
	return
}

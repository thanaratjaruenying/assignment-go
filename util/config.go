package util

import "github.com/spf13/viper"

type Config struct {
	Port string `mapstructure:"PORT"`
}

func LoadEnv() (*Config, error) {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	readErr := viper.ReadInConfig()
	if readErr != nil {
		return &Config{}, readErr
	}

	env := &Config{}

	err := viper.Unmarshal(env)
	if err != nil {
		return &Config{}, err
	}

	return env, nil
}

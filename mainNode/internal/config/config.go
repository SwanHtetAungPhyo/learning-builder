package config

import (
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/model"
	"github.com/spf13/viper"
)

type Server struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	Server     Server            `mapstructure:"server"`
	Validators []model.Validator `mapstructure:"validators"`
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

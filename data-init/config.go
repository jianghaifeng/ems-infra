package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

func getConfig(key string) interface{} {
	return viper.Get(key)
}

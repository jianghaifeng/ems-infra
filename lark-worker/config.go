package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetEnvPrefix("ems")
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	cfgFileName := "config"
	env := viper.GetString("env")
	if len(env) > 0 {
		cfgFileName += "-" + env
	}

	viper.SetConfigName(cfgFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func getConfigUInt(key string) uint64 {
	return viper.GetUint64(key)
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

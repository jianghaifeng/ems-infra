package main

import (
	"log"
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
	if err != nil {
		log.Panicf("fatal error config file: %v", err)
	}
}

func getConfigString(key string) string {
	return viper.GetString(key)
}

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Initiator() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	fmt.Println("Successfully read config file")
}

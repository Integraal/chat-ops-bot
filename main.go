package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
}
func main() {
	fmt.Println("Hi there")
}

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strconv"
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")
}
func main() {
	fmt.Println("Hi there")
}

func startBot() {
	bot, err := NewBot(TelegramConfig{})
	if err != nil {
		panic(err)
	}
	bot.OnAgree(func(chatId int64, eventId int64) {
		fmt.Println("User " + strconv.Itoa(int(chatId)) + " was present on event " + strconv.Itoa(int(eventId)))
	})
	bot.OnDisagree(func(chatId int64, eventId int64) {
		fmt.Println("User " + strconv.Itoa(int(chatId)) + " wasn't present on event " + strconv.Itoa(int(eventId)))
	})
	go bot.listen()
}
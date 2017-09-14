package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

var cal chan string
var users []user
var configuration config

type user struct {
	TelegramId   int
	JiraUsername string
	IcsLink      string
}
type config struct {
	Users []user
	Telegram TelegramConfig
}

func init() {
	conf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	json.Unmarshal(conf, &configuration)
}
func main() {

}

func startBot() {
	bot, err := NewBot(configuration.Telegram)
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
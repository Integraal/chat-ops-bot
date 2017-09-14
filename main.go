package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"sync"
	"github.com/Integraal/chat-ops-bot/telegram"
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
	Users []user `json:"users"`
	Telegram telegram.Config `json:"telegram"`
}

func init() {
	conf, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	json.Unmarshal(conf, &configuration)
}

func main() {
	var wg sync.WaitGroup
	bot := startBot(&wg)
	bot.SendPoll(999)
	wg.Wait()
}

func startBot(wg *sync.WaitGroup) *telegram.Bot {
	bot, err := telegram.NewBot(configuration.Telegram)
	if err != nil {
		panic(err)
	}
	bot.OnAgree(func(chatId int64, eventId int64) {
		fmt.Println("User " + strconv.Itoa(int(chatId)) + " was present on event " + strconv.Itoa(int(eventId)))
	})
	bot.OnDisagree(func(chatId int64, eventId int64) {
		fmt.Println("User " + strconv.Itoa(int(chatId)) + " wasn't present on event " + strconv.Itoa(int(eventId)))
	})
	go bot.Listen(wg)
	wg.Add(1)
	return bot
}
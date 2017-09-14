package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"sync"
	"github.com/Integraal/chat-ops-bot/telegram"
	"github.com/Integraal/chat-ops-bot/event"
)

var configuration config

type user struct {
	TelegramId   int
	JiraUsername string
	IcsLink      string
}

type config struct {
	Users    []user `json:"users"`
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
	event.Events = make(map[int64]event.Event)
	event.Events[1] = event.Event{
		ID: 1,
		Agreed: make(map[int64]bool),
		Disagreed: make(map[int64]bool),
	}
	event.Events[2] = event.Event{
		ID: 2,
		Agreed: make(map[int64]bool),
		Disagreed: make(map[int64]bool),
	}
	e := event.Events[2]
	var wg sync.WaitGroup
	bot := startBot(&wg)
	bot.SendPoll(&e)
	wg.Wait()
}

func startBot(wg *sync.WaitGroup) *telegram.Bot {
	bot, err := telegram.NewBot(configuration.Telegram)
	if err != nil {
		panic(err)
	}
	bot.OnAgree(func(chatId int64, eventId int64) *event.Event {
		if e, ok := event.Events[eventId]; ok {
			fmt.Println("NO")
			return &e
		}
		return nil
	})
	bot.OnDisagree(func(chatId int64, eventId int64) *event.Event  {
		if e, ok := event.Events[eventId]; ok {
			fmt.Println("YES")
			return &e
		}
		return nil
	})
	go bot.Listen(wg)
	wg.Add(1)
	return bot
}

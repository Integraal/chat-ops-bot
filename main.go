package main

import (
	"fmt"
	"sync"
	"github.com/integraal/chat-ops-bot/components/event"
	"github.com/integraal/chat-ops-bot/components/telegram"
	"github.com/integraal/chat-ops-bot/components/config"
	"github.com/integraal/chat-ops-bot/components/user"
)

var conf *config.Config

func init() {
	conf = config.Initialize()
	user.Initialize(conf.Users)
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
	bot, err := telegram.NewBot(conf.Telegram)
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

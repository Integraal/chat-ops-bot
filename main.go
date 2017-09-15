package main

import (
	"fmt"
	"strconv"
	"sync"
	"github.com/integraal/chat-ops-bot/telegram"
	"github.com/integraal/chat-ops-bot/components/config"
	"github.com/integraal/chat-ops-bot/components/user"
)

var conf *config.Config

func init() {
	conf = config.Initialize()
	user.Initialize(conf.Users)
}
func main() {
	var wg sync.WaitGroup
	bot := startBot(&wg)
	bot.SendPoll(999)
	wg.Wait()
}

func startBot(wg *sync.WaitGroup) *telegram.Bot {
	bot, err := telegram.NewBot(conf.Telegram)
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

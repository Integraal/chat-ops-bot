package main

import (
	"fmt"
	"sync"
	"github.com/integraal/chat-ops-bot/components/event"
	"github.com/integraal/chat-ops-bot/components/telegram"
	"github.com/integraal/chat-ops-bot/components/config"
	"github.com/integraal/chat-ops-bot/components/user"
	"github.com/integraal/chat-ops-bot/components/jira"
	"github.com/integraal/chat-ops-bot/components/watchdog"
)

var conf *config.Config

func init() {
	conf = config.Initialize()
	user.Initialize(conf.Users, 20)
	jira.Initialize(conf.Jira)
	watchdog.Initialize(conf.Watchdog)
}

func main() {
	var wg sync.WaitGroup
	startWatchdog(&wg)
	startBot(&wg)
	wg.Wait()
}

func startBot(wg *sync.WaitGroup) *telegram.Bot {
	bot, err := telegram.NewBot(conf.Telegram)
	if err != nil {
		panic(err)
	}
	bot.OnAgree(func(chatId int64, eventId string) *event.Event {
		e, err := event.Get(eventId)
		if err != nil {
			return nil
		}
		u, err := e.GetUser(chatId)
		if err != nil {
			return nil
		}
		issue, err := jira.Get().EnsureIssue(e)
		if err != nil {
			return nil
		}
		err = jira.Get().AddUserTime(issue, e, u)
		if err != nil {
			return nil
		}
		fmt.Println("YES")
		return e
	})
	bot.OnDisagree(func(chatId int64, eventId string) *event.Event {
		e, err := event.Get(eventId)
		if err != nil {
			return nil
		}
		fmt.Println("NO")
		return e
	})
	go bot.Listen(wg)
	wg.Add(1)
	return bot
}

func startWatchdog(wg *sync.WaitGroup) *watchdog.Watchdog {
	wd := watchdog.Get()

	wd.OnUpdate(fetchEvents) // Each X minutes
	wd.OnTick(func () {
		// Each minute
	})
	go wd.Listen(wg)
	wg.Add(1)
	return wd
}

func fetchEvents() {
	event.Clear()
	for _, u := range user.Get() {
		events, err := u.UpcomingEvents()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, e := range events {
			evt := event.NewEvent(&e)
			event.Append(evt, u)
		}
	}
}
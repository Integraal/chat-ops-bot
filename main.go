package main

import (
	"fmt"
	"sync"
	"github.com/integraal/chat-ops-bot/components/event"
	"github.com/integraal/chat-ops-bot/components/telegram"
	"github.com/integraal/chat-ops-bot/components/config"
	"github.com/integraal/chat-ops-bot/components/user"
	"github.com/integraal/chat-ops-bot/components/jira"
)

var conf *config.Config

func init() {
	conf = config.Initialize()
	user.Initialize(conf.Users, 20)
	jira.Initialize(conf.Jira)
}

func main() {
	fetchEvents()
	fmt.Println(event.GetAll())
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
	//evt, err := event.Get("4f3aab81210fe3ec020e10ce77b21e57")
	//if err != nil {
	//	panic(err)
	//}
	//issue, err := jira.Get().EnsureIssue(evt)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%v", issue)
	//var wg sync.WaitGroup
	//bot := startBot(&wg)
	//bot.SendPoll(evt)
	//wg.Wait()
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

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
	"github.com/integraal/chat-ops-bot/components/datebook"
	"time"
)

var conf *config.Config

func init() {
	conf = config.Initialize()
	user.Initialize(conf.Users)
	jira.Initialize(conf.Jira)
	datebook.Initialize(conf.Calendar.UpcomingLimit)
	watchdog.Initialize(conf.Watchdog)
}

func main() {
	var wg sync.WaitGroup
	bot := startBot(&wg)
	startWatchdog(&wg, bot)
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
			fmt.Println(err)
			return nil
		}
		return e
	})
	bot.OnDisagree(func(chatId int64, eventId string) *event.Event {
		e, err := event.Get(eventId)
		if err != nil {
			return nil
		}
		return e
	})
	go bot.Listen(wg)
	wg.Add(1)
	return bot
}

func startWatchdog(wg *sync.WaitGroup, bot *telegram.Bot) *watchdog.Watchdog {
	wd := watchdog.Get()

	wd.OnUpdate(fetchEvents) // Each X minutes
	wd.OnTick(func() {
		fmt.Println("Tick.")
		events := event.GetAll()
		now := time.Now()
		fmt.Println(time.Now().Format("02.01.2006 15:04:05 -0700"))
		for eventId := range *events {
			e, _ := event.Get(eventId)
			// Check if event is upcoming
			toStart := e.Start.Sub(now)
			remind := toStart > 0
			remind = remind && toStart <= time.Duration(wd.RemindBefore)*time.Minute
			if remind && !e.ReminderSent {
				e.ReminderSent = true
				bot.SendReminder(e)
			}
			// Check if event finished
			fromEnd := now.Sub(e.End)
			sendPoll := fromEnd >= time.Duration(wd.RemindAfter)*time.Minute
			sendPoll = sendPoll && fromEnd <= time.Duration(wd.DontRemindAfter)*time.Minute
			fmt.Println(e.Summary, e.End.Format("02.01.2006 15:04:05 -0700"), fromEnd)
			if sendPoll && !e.PollSent {
				e.PollSent = true
				bot.SendPoll(e)
			}
		}
	})
	go wd.Listen(wg)
	wg.Add(1)
	return wd
}

func fetchEvents() {
	fmt.Println("Fetching new events...")
	for _, u := range user.Get() {
		events, err := u.DatesAround(time.Now(), 2, 2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, e := range events {
			evt := event.NewEvent(e)
			event.Append(&evt, u)
		}
	}
	fmt.Println("Done.")
}

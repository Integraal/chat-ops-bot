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
	"github.com/integraal/chat-ops-bot/components/db"
	"io/ioutil"
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
	bot.OnAgree(func(chatId int64, eventId string) (*event.Event, error) {
		e, err := event.Get(eventId)
		if err != nil {
			return nil, err
		}
		u, err := e.GetUser(chatId)
		if err != nil {
			return nil, err
		}
		issue, response, err := jira.Get().EnsureIssue(e)
		if err != nil {
			if response != nil {
				body, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(body))
			}
			return nil, err
		}
		response, err = jira.Get().AddUserTime(issue, e, u)
		if err != nil {
			if response != nil {
				body, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(body))
			}
			return nil, err
		}
		return e, nil
	})
	bot.OnDisagree(func(chatId int64, eventId string) (*event.Event, error) {
		e, err := event.Get(eventId)
		if err != nil {
			return nil, err
		}
		return e, nil
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
			evt, _ := event.Get(eventId)
			dbEvent := db.Get().Event(eventId)
			// Check if event is upcoming
			toStart := evt.Start.Sub(now)
			remind := toStart > 0
			remind = remind && toStart <= time.Duration(wd.RemindBefore)*time.Minute
			if remind && !dbEvent.GetReminderSent() {
				dbEvent.SetReminderSent(true)
				bot.SendReminder(evt)
			}
			// Check if event finished
			fromEnd := now.Sub(evt.End)
			sendPoll := fromEnd >= time.Duration(wd.RemindAfter)*time.Minute
			sendPoll = sendPoll && fromEnd <= time.Duration(wd.DontRemindAfter)*time.Minute
			fmt.Println(evt.Summary, evt.End.Format("02.01.2006 15:04:05 -0700"), fromEnd)
			if sendPoll && !dbEvent.GetPollSent() {
				dbEvent.SetPollSent(true)
				bot.SendPoll(evt)
			}
		}
	})
	go wd.Listen(wg)
	wg.Add(1)
	return wd
}

func fetchEvents() {
	fmt.Println("Fetching new events...")
	event.Clear()
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

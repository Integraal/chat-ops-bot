package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"encoding/json"
)

type TelegramConfig struct {
	Token   string `json:"token"`
	ChatID  int64 `json:"chatId"`
	Timeout int `json:"timeout"`
}

var bot Bot

type Bot struct {
	timeout    int
	chatId     int64
	botApi     telegram.BotAPI
	onAgree    func(chatId int64, eventId int64)
	onDisagree func(chatId int64, eventId int64)
}

const (
	REPLY_YES = "Yes"
	REPLY_NO = "No"
)

type ButtonPress struct {
	eventId int64 // TODO: event struct
	reply string
}

func (b *Bot) OnAgree(callback func(chatId int64, eventId int64)) {
	b.onAgree = callback
}

func (b *Bot) OnDisagree(callback func(chatId int64, eventId int64)) {
	b.onDisagree = callback
}

func NewBot(config TelegramConfig) (Bot, error) {
	botApi, err := telegram.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot = Bot{
		chatId: config.ChatID,
		timeout: config.Timeout,
		botApi:  *botApi,
	}
	return bot, nil
}

func (b *Bot) listen() error {
	if b.onAgree == nil {
		return error("Bot onAgree callback is not set")
	}
	if b.onDisagree == nil {
		return error("Bot onDisagree callback is not set")
	}
	u := telegram.NewUpdate(0)
	u.Timeout = b.timeout

	updates, err := b.botApi.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			buttonPress := ButtonPress{}
			err = json.Unmarshal([]byte(update.CallbackQuery.Data), buttonPress)
			if err != nil {
				panic(err)
			}
			if buttonPress.reply == REPLY_YES {
				b.onAgree(int64(update.CallbackQuery.From.ID), buttonPress.eventId)
			}
			if buttonPress.reply == REPLY_NO {
				b.onDisagree(int64(update.CallbackQuery.From.ID), buttonPress.eventId)
			}
		}
	}
	return nil
}

//TODO: Event struct and pretty printed message
func (b *Bot) SendReminder(eventId int64) {
	text := "Will you go to event" + strconv.Itoa(int(eventId)) + "?"
	message := telegram.NewMessage(b.chatId, text)
	b.botApi.Send(message)
}

func (bp *ButtonPress) marshall() string {
	str, err := json.Marshal(bp)
	if err != nil {
		panic(err)
	}
	return string(str)
}

//TODO: Event struct and pretty printed message
func (b *Bot) SendPoll(eventId int64) {
	text := "Did you get to event" + strconv.Itoa(int(eventId)) + "?"
	message := telegram.NewMessage(b.chatId, text)

	yes := ButtonPress{
		eventId: eventId,
		reply: REPLY_YES,
	}
	no := ButtonPress{
		eventId: eventId,
		reply: REPLY_NO,
	}

	buttons := []telegram.InlineKeyboardButton{
		telegram.NewInlineKeyboardButtonData("Yes", yes.marshall()),
		telegram.NewInlineKeyboardButtonData("No", no.marshall()),
	}

	message.ReplyMarkup = telegram.NewInlineKeyboardMarkup(buttons)
	b.botApi.Send(message)
}

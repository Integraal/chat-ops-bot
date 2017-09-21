package telegram

import (
	tlg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"encoding/json"
	"sync"
	"github.com/integraal/chat-ops-bot/components/event"
	"fmt"
)

type Config struct {
	Token   string `json:"token"`
	ChatID  int64 `json:"chatId"`
	Timeout int `json:"timeout"`
}

const (
	REPLY_YES = "Yes"
	REPLY_NO  = "No"
)

type ButtonPress struct {
	EventID string `json:"eventId"`
	Reply   string `json:"reply"`
}

func (bp *ButtonPress) marshall() string {
	str, err := json.Marshal(bp)
	if err != nil {
		panic(err)
	}
	return string(str)
}

var bot Bot

type CallbackFunc func(chatId int64, eventId string) (*event.Event, error)

type Bot struct {
	timeout    int
	chatId     int64
	botApi     tlg.BotAPI
	onAgree    CallbackFunc
	onDisagree CallbackFunc
}

func (b *Bot) OnAgree(callback CallbackFunc) {
	b.onAgree = callback
}

func (b *Bot) OnDisagree(callback CallbackFunc) {
	b.onDisagree = callback
}

func NewBot(config Config) (*Bot, error) {
	botApi, err := tlg.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}
	bot = Bot{
		chatId:  config.ChatID,
		timeout: config.Timeout,
		botApi:  *botApi,
	}
	return &bot, nil
}

func (b *Bot) Listen(wg *sync.WaitGroup) {
	if b.onAgree == nil {
		panic("Bot onAgree callback is not set")
	}
	if b.onDisagree == nil {
		panic("Bot onDisagree callback is not set")
	}
	b.botApi.Debug = true
	u := tlg.NewUpdate(0)
	u.Timeout = b.timeout

	updates, err := b.botApi.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			var buttonPress ButtonPress
			err = json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonPress)
			if err != nil {
				panic(err)
			}
			var responseEvent *event.Event
			userId := int64(update.CallbackQuery.From.ID)
			if buttonPress.Reply == REPLY_YES {
				responseEvent, err = b.onAgree(userId, buttonPress.EventID)
				if err == nil {
					responseEvent.SetAttended(userId, true)
					b.botApi.AnswerCallbackQuery(tlg.NewCallback(update.CallbackQuery.ID, "Ок, я отмечу время в JIRA"))
				} else {
					fmt.Println(err)
					b.botApi.AnswerCallbackQuery(tlg.NewCallback(update.CallbackQuery.ID, "Извините, произошла какая-то ошибка"))
				}
			}
			if buttonPress.Reply == REPLY_NO {
				responseEvent, err = b.onDisagree(userId, buttonPress.EventID)
				if err == nil {
					responseEvent.SetAttended(userId, false)
					b.botApi.AnswerCallbackQuery(tlg.NewCallback(update.CallbackQuery.ID, "Время в JIRA не будет учтено"))
				} else {
					fmt.Println(err)
					b.botApi.AnswerCallbackQuery(tlg.NewCallback(update.CallbackQuery.ID, "Извините, произошла какая-то ошибка"))
				}
			}
			if responseEvent != nil {
				b.updatePollMarkup(responseEvent, update.CallbackQuery.Message.MessageID)
			}
		}
	}
	wg.Done()
}

func (b *Bot) getPollMarkup(event *event.Event) tlg.InlineKeyboardMarkup {

	yes := ButtonPress{
		EventID: event.Id,
		Reply:   REPLY_YES,
	}
	no := ButtonPress{
		EventID: event.Id,
		Reply:   REPLY_NO,
	}

	yesText := "👍"
	noText := "👎"

	if count := event.GetAttendedCount(); count > 0 {
		yesText += " " + strconv.Itoa(count)
	}
	if count := event.GetUnattendedCount(); count > 0 {
		noText += " " + strconv.Itoa(count)
	}

	buttons := []tlg.InlineKeyboardButton{
		tlg.NewInlineKeyboardButtonData(yesText, yes.marshall()),
		tlg.NewInlineKeyboardButtonData(noText, no.marshall()),
	}
	return tlg.NewInlineKeyboardMarkup(buttons)
}

func (b *Bot) updatePollMarkup(event *event.Event, messageId int) {
	message := tlg.NewEditMessageReplyMarkup(b.chatId, messageId, b.getPollMarkup(event))
	b.botApi.Send(message)
}

func (b *Bot) SendPoll(event *event.Event) {
	text := `*%s*
Кто присутствовал на данном мероприятии?
*%s* - *%s*, *%s*
`
	msg := fmt.Sprintf(
		text,
		event.Summary,
		event.GetStartTime().Format("15:04"),
		event.GetEndTime().Format("15:04"),
		event.GetStartTime().Format("02.01.2006"),
	)
	message := tlg.NewMessage(b.chatId, msg)
	message.ParseMode = "Markdown"
	message.ReplyMarkup = b.getPollMarkup(event)
	b.botApi.Send(message)
}

func (b *Bot) SendReminder(event *event.Event) {
	text := `
Приготовьтесь ко встрече *%s*
*%s* - *%s*, *%s*
Участники:
`
	text = fmt.Sprintf(
		text,
		event.Summary,
		event.GetStartTime().Format("15:04"),
		event.GetEndTime().Format("15:04"),
		event.GetStartTime().Format("02.01.2006"),
	)
	for _, user := range event.GetUsers() {
		text += fmt.Sprintf(" • [%s](tg://user?id=%s)", user.Name, strconv.Itoa(user.TelegramId))
	}
	b.sendMessage(text)
}

func (b *Bot) sendMessage(text string) {
	message := tlg.NewMessage(b.chatId, text)
	message.ParseMode = "Markdown"
	b.botApi.Send(message)
}

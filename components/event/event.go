package event

import (
	ics "github.com/PuloV/ics-golang"
	"github.com/integraal/chat-ops-bot/components/user"
	"errors"
	"time"
)

var events map[string]Event = make(map[string]Event)

type Event struct {
	ID          string
	Summary     string
	Description string
	Duration time.Duration

	Users map[int64]user.User
	Agreed map[int64]bool
	Disagreed map[int64]bool
}

func Clear() {
	events = make(map[string]Event)
}

func (e *Event) SetAgree(userId int64) {
	e.Agreed[userId] = true
	if _, ok := e.Disagreed[userId]; ok {
		delete(e.Disagreed, userId)
	}
}
func (e *Event) SetDisagree(userId int64) {
	e.Disagreed[userId] = true
	if _, ok := e.Agreed[userId]; ok {
		delete(e.Agreed, userId)
	}
}

func Append(event Event, user user.User) {
	if e, ok := events[event.ID]; ok {
		events[e.ID].Users[int64(user.ChatId)] = user
	} else {
		event.Users[int64(user.ChatId)] = user
		events[event.ID] = event
	}
}

func NewEvent(ics *ics.Event) Event {
	evt := Event{
		ID:      ics.GetID(),
		Summary: ics.GetSummary(),
		Description: ics.GetDescription(),
		Duration: ics.GetEnd().Sub(ics.GetStart()),

		Users: make(map[int64]user.User),
		Agreed: make(map[int64]bool),
		Disagreed: make(map[int64]bool),
	}
	return evt
}

func GetAll() *map[string]Event {
	return &events
}

func Get(eventId string) (*Event, error) {
	if e, ok := events[eventId]; ok {
		return &e, nil
	}
	return nil, errors.New("Event does not exists")
}

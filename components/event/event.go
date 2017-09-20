package event

import (
	"github.com/integraal/ics-golang"
	"github.com/integraal/chat-ops-bot/components/user"
	"errors"
	"time"
	"strings"
	"github.com/integraal/chat-ops-bot/components/db"
)

var events map[string]*Event = make(map[string]*Event)

type Event struct {
	ID          string
	ImportedID  string
	Summary     string
	Description string
	Duration    time.Duration
	Start       time.Time
	End         time.Time

	users map[int64]user.User

	dbEvent		*db.Event
}

func Clear() {
	events = make(map[string]*Event)
}

func (e *Event) GetUser(chatId int64) (*user.User, error) {
	if u, ok := e.users[chatId]; ok {
		return &u, nil
	}
	return nil, errors.New("User does not exist")
}

func (e *Event) GetUsers() *map[int64]user.User {
	return &e.users
}

func Append(event *Event, user user.User) {
	if e, ok := events[event.ID]; ok {
		events[e.ID].users[int64(user.TelegramId)] = user
	} else {
		event.users[int64(user.TelegramId)] = user
		events[event.ID] = event
	}
}

func NewEvent(ics *ics.Event) Event {
	evt := Event{
		ID:          ics.GetID(),
		ImportedID:  strings.Replace(ics.GetImportedID(), "-", "", -1),
		Summary:     ics.GetSummary(),
		Description: ics.GetDescription(),
		Duration:    ics.GetEnd().Sub(ics.GetStart()),
		Start:       ics.GetStart(),
		End:         ics.GetEnd(),

		users: make(map[int64]user.User),
	}
	return evt
}

func GetAll() *map[string]*Event {
	return &events
}

func Get(id string) (*Event, error) {
	if e, ok := events[id]; ok {
		return e, nil
	}
	return nil, errors.New("Event does not exist")
}

func (e *Event) getDbEvent() *db.Event {
	if e.dbEvent == nil {
		e.dbEvent = db.Get().Event(e.ID)
	}
	return e.dbEvent
}

func (e *Event) SetAttended(userId int64, value bool) {
	e.getDbEvent().SetAttended(userId, value)
}

func (e *Event) GetAttendedCount() int {
	return e.getDbEvent().GetAttendedCount()
}

func (e *Event) GetUnattendedCount() int {
	return e.getDbEvent().GetUnattendedCount()
}

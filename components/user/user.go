package user

import (
	"github.com/PuloV/ics-golang"
	"github.com/integraal/chat-ops-bot/components/datebook"
)

var usersArray []User

type User struct {
	TelegramId   int
	JiraUsername string
	IcsLinks     []string
	Events       []*ics.Event
}

func Initialize(users []User) {
	usersArray = users
}
func Get() []User {
	return usersArray
}
func (u *User) Calendars() ([]*ics.Calendar, error) {
	return datebook.Calendars(u.IcsLinks)
}

func (u *User) UpcomingEvents() ([]*ics.Event, error) {
	var events []*ics.Event
	calendars, err := u.Calendars()
	if err != nil {
		return nil, err
	}
	for _, calendar := range calendars {
		events = append(events, datebook.UpcomingEvents(calendar)...)
	}
	u.Events = events
	return events, nil
}

func UpcomingEvents() ([]*ics.Event, error) {
	var events []*ics.Event
	for _, user := range Get() {
		evs, err := user.UpcomingEvents()
		if err != nil {
			return nil, err
		}
		events = append(events, evs...)
	}
	return events, nil
}

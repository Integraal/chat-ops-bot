package user

import (
	"github.com/PuloV/ics-golang"
	"time"
)

var usersArray []User

type User struct {
	ChatId       int
	JiraUsername string
	IcsLink      string
}

type CalendarConfig struct {
	UpcomingEvents int
}

var upcomingLimit int

func Initialize(users []User, limit int) {
	usersArray = users
	upcomingLimit = limit
}

func Get() []User {
	return usersArray
}
func (u *User) Calendars() ([]*ics.Calendar, error) {
	parser := ics.New()
	inputChan := parser.GetInputChan()
	inputChan <- u.IcsLink
	parser.Wait()
	return parser.GetCalendars()
}

func (u *User) Events() ([]ics.Event, error) {
	var events []ics.Event
	cals, err := u.Calendars()
	if err != nil {
		return nil, err
	}
	for _, cal := range cals {
		events = append(events, cal.GetEvents()...)
	}
	return events, nil
}
func (u *User) UpcomingEvents() ([]ics.Event, error) {
	var events []ics.Event
	cals, err := u.Calendars()
	if err != nil {
		return nil, err
	}
	for _, cal := range cals {
		events = append(events, cal.GetUpcomingEvents(upcomingLimit)...)
	}
	return events, nil
}
func (u *User) EventsByTime(t time.Time) ([]*ics.Event, error) {
	var events []*ics.Event
	cals, err := u.Calendars()
	if err != nil {
		return nil, err
	}
	for _, cal := range cals {
		evs, err := cal.GetEventsByDate(t)
		if err != nil {
			return nil, err
		}
		events = append(events, evs...)
	}
	return events, nil
}
func Calendars() ([]*ics.Calendar, error) {
	var calendars []*ics.Calendar
	for _, user := range Get() {
		cal, err := user.Calendars()
		if err != nil {
			return nil, err
		}
		calendars = append(calendars, cal...)
	}
	return calendars, nil
}
func UpcomingEvents() ([]ics.Event, error) {
	var events []ics.Event
	for _, user := range Get() {
		evs, err := user.UpcomingEvents()
		if err != nil {
			return nil, err
		}
		events = append(events, evs...)
	}
	return events, nil
}

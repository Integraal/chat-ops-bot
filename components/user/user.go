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

func Initialize(users []User) {
	usersArray = users
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

func (u *User) Events(t time.Time) ([]*ics.Event, error) {
	cal, err := u.Calendars()
	if err != nil {
		return nil, err
	}
	return cal[0].GetEventsByDate(t)
}

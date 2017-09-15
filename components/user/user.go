package user

import (
	"github.com/PuloV/ics-golang"
	"time"
)

var usersArray []User

type User struct {
	TelegramId   int
	JiraUsername string
	IcsLink      string
}

func Initialize(users []User) {
	usersArray = users
}
func Get() []User {
	return usersArray
}
func (u *User) Events(t time.Time) ([]*ics.Event, error) {
	parser := ics.New()
	inputChan := parser.GetInputChan()
	inputChan <- u.IcsLink
	parser.Wait()
	cal, err := parser.GetCalendars()
	if err != nil {
		return nil, err
	}
	return cal[0].GetEventsByDate(t)
}

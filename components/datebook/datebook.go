package datebook

import (
	"github.com/PuloV/ics-golang"
)

type CalendarConfig struct {
	UpcomingLimit int
}

func Initialize(limit int) {
	upcomingLimit = limit
}

var upcomingLimit int

func UpcomingEvents(calendar *ics.Calendar) []*ics.Event {
	var events []*ics.Event
	for _, event := range calendar.GetUpcomingEvents(upcomingLimit) {
		e := event
		events = append(events, &e)
	}
	return events
}
func Calendars(links []string) ([]*ics.Calendar, error) {
	parser := ics.New()
	inputChan := parser.GetInputChan()
	for _, link := range links {
		inputChan <- link
	}
	parser.Wait()
	return parser.GetCalendars()
}

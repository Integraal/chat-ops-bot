package datebook

import (
	"github.com/integraal/ics-golang"
	"time"
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
	ics.RepeatRuleApply = true
	inputChan := parser.GetInputChan()
	for _, link := range links {
		inputChan <- link
	}
	parser.Wait()
	return parser.GetCalendars()
}
func DatesAround(calendar *ics.Calendar, date time.Time, daysBefore, daysAfter int) ([]*ics.Event, error) {
	var events []*ics.Event
	tBefore := date.AddDate(0, 0, -daysBefore)
	tAfter := date.AddDate(0, 0, daysAfter+1)
	for i := tBefore; i != tAfter; i = i.AddDate(0, 0, 1) {
		dayEvents, _ := calendar.GetEventsByDate(i)
		for _, event := range dayEvents {
			events = append(events, event)
		}
	}
	return events, nil
}

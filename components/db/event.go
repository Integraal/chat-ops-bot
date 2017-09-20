package db

type Event struct {
	id string
	reminderSent bool
	pollSent bool
	attended map[int64]bool
}

func (d *Database) Event(eventId string) *Event {
	if _, ok := database.events[eventId]; !ok {
		d.events[eventId] = &Event{
			id: eventId,
			reminderSent: false,
			pollSent:     false,
			attended:     make(map[int64]bool),
		}
	}
	return d.events[eventId]
}

func (e *Event) SetAttended(userId int64, attended bool) {
	e.attended[userId] = attended
}

func (e *Event) GetAttendedCount() int {
	count := 0
	for _, attended := range e.attended {
		if attended == true {
			count++
		}
	}
	return count
}

func (e *Event) GetUnattendedCount() int {
	count := 0
	for _, attended := range e.attended {
		if attended == false {
			count++
		}
	}
	return count
}

func (e *Event) SetReminderSent(value bool) {
	e.reminderSent = value
}

func (e *Event) GetReminderSent() bool {
	return e.reminderSent
}

func (e *Event) SetPollSent(value bool) {
	e.pollSent = value
}

func (e *Event) GetPollSent() bool {
	return e.pollSent
}

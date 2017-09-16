package event

var Events map[int64]Event

type Event struct {
	ID int64
	Name string
	Description string

	Agreed map[int64]bool
	Disagreed map[int64]bool
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
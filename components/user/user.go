package user

var usersArray []User

type User struct {
	Name         string
	Email        string
	TelegramId   int
	JiraUsername string
}

func Initialize(users []User) {
	usersArray = users
}

func Get() []User {
	return usersArray
}

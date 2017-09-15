package user

var usersArray []User

type User struct {
	TelegramId   int
	JiraUsername string
	IcsLink      string
}
func Initialize(users []User){
	usersArray = users
}
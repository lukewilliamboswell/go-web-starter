package user

type UserRepository interface {
	InsertUser(User) error
	GetUsers() ([]User, error)
	GetUser(string) (User, error)
}

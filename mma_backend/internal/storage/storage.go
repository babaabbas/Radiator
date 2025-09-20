package storage

type Storage interface {
	Create_User() error
	User_Login() error
}

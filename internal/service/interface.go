package service

type UserService interface {
	GetUser()
	CreateUser()
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}

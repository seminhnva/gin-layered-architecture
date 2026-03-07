package repository

type UserRepository interface {
	FindUser()
	Create()
	FindByUUID()
	Update()
	Delete()
}

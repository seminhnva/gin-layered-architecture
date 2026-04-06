package auth

type Hasher interface {
	HashPassword(plain_password string) (string, error)
	CheckPasswordHash(plain_password, hashed string) (bool, error)
}

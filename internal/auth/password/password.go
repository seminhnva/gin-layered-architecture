package password

import (
	"github.com/alexedwards/argon2id"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
)

type Argon2Hasher struct {
}

func NewPasswordService() auth.Hasher {
	return &Argon2Hasher{}
}

func (ps *Argon2Hasher) HashPassword(plain_password string) (string, error) {
	hash, err := argon2id.CreateHash(plain_password, argon2id.DefaultParams)
	return hash, err
}

func (ps *Argon2Hasher) CheckPasswordHash(plain_password, hashed string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(plain_password, hashed)
	return match, err
}

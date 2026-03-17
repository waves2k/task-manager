package service

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	GeneratePasswordHash(password string) string
	VerifyPassword(password, passwordHash string) bool
}

type passwordHasher struct {
	salt string
}

func NewPasswordHasher(passwordSalt string) *passwordHasher {
	return &passwordHasher{
		salt: passwordSalt,
	}
}

func (h *passwordHasher) GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}

func (h *passwordHasher) VerifyPassword(password, passwordHash string) bool {
	return h.GeneratePasswordHash(password) == passwordHash
}

package common

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	Hash(password string) (string, error)
	CheckHash(hashedPassword string, password string) bool
}

type PasswordHash struct {
	token string
}

func NewPasswordHash(token string) PasswordHasher {
	return &PasswordHash{
		token: token,
	}
}

func (hash PasswordHash) getPattern(password string) []byte {
	return []byte(password + hash.token)
}

func (hash PasswordHash) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(hash.getPattern(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (hash PasswordHash) CheckHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), hash.getPattern(password))
	return err == nil
}

package passwordhash

import "golang.org/x/crypto/bcrypt"

type PasswordHash interface {
	Hash(password string) (string, error)
	CheckHash(hashedPassword string, password string) bool
}

type passwordHash struct {
	token string
}

func New(token string) PasswordHash {
	return &passwordHash{
		token: token,
	}
}

func (hash passwordHash) getPattern(password string) []byte {
	return []byte(password + hash.token)
}

func (hash passwordHash) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(hash.getPattern(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func (hash passwordHash) CheckHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), hash.getPattern(password))
	return err == nil
}

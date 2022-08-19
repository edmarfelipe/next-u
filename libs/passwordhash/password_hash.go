package passwordhash

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/tracer"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHash interface {
	Hash(ctx context.Context, password string) (string, error)
	CheckHash(ctx context.Context, hashedPassword string, password string) bool
}

type passwordHash struct {
	token  string
	logger logger.Logger
}

func New(token string, logger logger.Logger) PasswordHash {
	return &passwordHash{
		token:  token,
		logger: logger,
	}
}

func (hash passwordHash) getPattern(password string) []byte {
	return []byte(password + hash.token)
}

func (hash passwordHash) Hash(ctx context.Context, password string) (string, error) {
	_, span := tracer.StartSpan(ctx, "PasswordHash", "Hash")
	defer span.End()

	hash.logger.Info(ctx, "Hashing password")

	bytes, err := bcrypt.GenerateFromPassword(hash.getPattern(password), 14)
	if err != nil {
		hash.logger.Error(ctx, "Fail to hash password", err)
		return "", err
	}

	hash.logger.Error(ctx, "Password hashed", err)
	return string(bytes), err
}

func (hash passwordHash) CheckHash(ctx context.Context, hashedPassword string, password string) bool {
	_, span := tracer.StartSpan(ctx, "PasswordHash", "CheckHash")
	defer span.End()

	hash.logger.Info(ctx, "Checking password hash")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), hash.getPattern(password))
	return err == nil
}

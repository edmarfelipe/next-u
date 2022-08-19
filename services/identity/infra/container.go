package infra

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/mail"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/libs/tracer"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Container struct {
	Config       *Config
	Logger       logger.Logger
	MailService  mail.MailService
	UserDB       db.UserDB
	PasswordHash passwordhash.PasswordHash
}

func NewContainer() *Container {
	logger := logger.New()
	ctx := context.Background()

	err := tracer.NewTracer()
	if err != nil {
		logger.Error(ctx, err.Error())
	}

	config, err := NewConfig("./config.yaml")
	if err != nil {
		logger.Error(ctx, err.Error())
	}

	mongoDB, err := db.NewConnection(config.DataBase.Name, config.DataBase.URI)
	if err != nil {
		logger.Error(ctx, err.Error())
	}

	return &Container{
		Logger: logger,
		Config: config,
		MailService: mail.New(
			logger.With("mail"),
			mail.ConfigEmail{
				Title:  config.Title,
				Email:  config.Email,
				ApiKey: config.SendGrid.APIKey,
			},
		),
		UserDB:       db.NewUser(mongoDB, logger.With("db")),
		PasswordHash: passwordhash.New(config.PasswordToken, logger.With("password-hash")),
	}
}

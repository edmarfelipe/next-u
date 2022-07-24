package infra

import (
	"log"

	"github.com/edmarfelipe/next-u/libs/mail"
	"github.com/edmarfelipe/next-u/libs/passwordhash"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Container struct {
	Config          *Config
	MailService     mail.MailService
	Validator       Validatorer
	UserDB          db.UserDB
	PasswordResetDB db.PasswordResetDB
	PasswordHash    passwordhash.PasswordHash
}

func NewContainer() *Container {
	config, err := NewConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	mongoDB, err := db.NewConnection(config.DataBase.Name, config.DataBase.URI)
	if err != nil {
		log.Fatal(err)
	}

	return &Container{
		Config:          config,
		MailService:     mail.New(config.Title, config.Email, config.SendGrid.APIKey),
		Validator:       NewValidator(),
		UserDB:          db.NewUser(mongoDB),
		PasswordResetDB: db.NewPasswordReset(mongoDB),
		PasswordHash:    passwordhash.New(config.PasswordToken),
	}
}

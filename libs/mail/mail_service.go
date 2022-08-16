package mail

import (
	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ConfigEmail struct {
	Title  string
	Email  string
	ApiKey string
}

type MailService interface {
	Send(mailTo MailTo, subject string, content string) error
}

type mailService struct {
	logger logger.Logger
	config ConfigEmail
}

func New(logger logger.Logger, config ConfigEmail) MailService {
	return &mailService{
		logger: logger,
		config: config,
	}
}

func (m *mailService) Send(mailTo MailTo, subject string, content string) error {
	from := mail.NewEmail(m.config.Title, m.config.Email)
	to := mail.NewEmail(mailTo.Name, mailTo.Email)

	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(m.config.ApiKey)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

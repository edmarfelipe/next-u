package infra

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailServicer interface {
	Send(mailTo MailTo, subject string, content string) error
}

type MailService struct {
	config Config
}

func NewMail(config Config) MailServicer {
	return &MailService{
		config: config,
	}
}

func (obj *MailService) Send(mailTo MailTo, subject string, content string) error {
	from := mail.NewEmail(obj.config.Title, obj.config.Email)
	to := mail.NewEmail(mailTo.Name, mailTo.Email)

	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(obj.config.APIKey)
	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

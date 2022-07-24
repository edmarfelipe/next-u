package mail

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type MailService interface {
	Send(mailTo MailTo, subject string, content string) error
}

type mailService struct {
	title  string
	email  string
	apiKey string
}

func New(title string, email string, apiKey string) MailService {
	return &mailService{
		title:  title,
		email:  email,
		apiKey: apiKey,
	}
}

func (m *mailService) Send(mailTo MailTo, subject string, content string) error {
	from := mail.NewEmail(m.title, m.email)
	to := mail.NewEmail(mailTo.Name, mailTo.Email)

	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(m.apiKey)

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

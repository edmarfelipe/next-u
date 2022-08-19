package mail

import (
	"context"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/tracer"
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
	Send(ctx context.Context, mailTo MailTo, subject string, content string) error
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

func (m *mailService) Send(ctx context.Context, mailTo MailTo, subject string, content string) error {
	_, span := tracer.StartSpan(ctx, "MailService", "Send")
	defer span.End()

	m.logger.Info(ctx, "Sending mail to "+mailTo.Email)

	from := mail.NewEmail(m.config.Title, m.config.Email)
	to := mail.NewEmail(mailTo.Name, mailTo.Email)

	message := mail.NewSingleEmail(from, subject, to, "", content)
	client := sendgrid.NewSendClient(m.config.ApiKey)

	_, err := client.Send(message)
	if err != nil {
		m.logger.Info(ctx, "Fail to send mail", err)
		return err
	}

	m.logger.Info(ctx, "Mail sent to "+mailTo.Email)

	return nil
}

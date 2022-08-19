package reset

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/edmarfelipe/next-u/libs/logger"
	"github.com/edmarfelipe/next-u/libs/mail"
	"github.com/edmarfelipe/next-u/libs/validator"
	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

var (
	errUserNotFound = errors.New("could not found user")
)

type Usecase interface {
	Execute(ctx context.Context, input Input) error
}

type Input struct {
	Email string `json:"email" validate:"email"`
}

type usecase struct {
	logger      logger.Logger
	config      *infra.Config
	userDB      db.UserDB
	mailService mail.MailService
}

func NewUsecase(
	logger logger.Logger,
	config *infra.Config,
	userDB db.UserDB,
	mailService mail.MailService,
) Usecase {
	return &usecase{
		logger:      logger,
		config:      config,
		userDB:      userDB,
		mailService: mailService,
	}
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	usc.logger.Info(ctx, "Recovering password from user "+in.Email)
	err := validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	if user == nil {
		usc.logger.Error(ctx, "Failed to recovery", "err", errUserNotFound)
		return errUserNotFound
	}

	mailTo := mail.MailTo{
		Name:  user.Name,
		Email: user.Email,
	}

	usc.logger.Info(ctx, "Creating password token")
	reset := user.CreatePasswordToken()
	err = usc.userDB.Update(ctx, *user)
	if err != nil {
		return err
	}

	subject := usc.config.Title + " - Reset Password"
	usc.mailService.Send(ctx, mailTo, subject, usc.createMailContent(ctx, reset.Token))
	if err != nil {
		return err
	}

	return nil
}

func (usc *usecase) createMailContent(ctx context.Context, token string) string {
	usc.logger.Info(ctx, "Building mail content")
	template := `
		<h1>%s</h1>
		<br/>
		<a href="%s">Reset Password</a>
	`

	changePasswordUrl := strings.ReplaceAll(
		usc.config.UrlPageChangePassword,
		"{token}",
		token,
	)

	return fmt.Sprintf(template, usc.config.Title, changePasswordUrl)
}

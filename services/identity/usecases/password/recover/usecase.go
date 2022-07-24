package recover

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/edmarfelipe/next-u/libs/mail"
	"github.com/edmarfelipe/next-u/services/identity/entity"
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
	config          *infra.Config
	userDB          db.UserDB
	passwordResetDB db.PasswordResetDB
	mailService     mail.MailService
	validator       infra.Validatorer
}

func NewUsecase(
	config *infra.Config,
	userRepository db.UserDB,
	passwordResetDB db.PasswordResetDB,
	mailService mail.MailService,
	validator infra.Validatorer,
) Usecase {
	return &usecase{
		config:          config,
		userDB:          userRepository,
		passwordResetDB: passwordResetDB,
		mailService:     mailService,
		validator:       validator,
	}
}

func (usc *usecase) Execute(ctx context.Context, in Input) error {
	err := usc.validator.IsValid(in)
	if err != nil {
		return err
	}

	user, err := usc.userDB.FindByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	if user == nil {
		return errUserNotFound
	}

	mailTo := mail.MailTo{
		Name:  user.Name,
		Email: user.Email,
	}

	passwordRest := entity.MakePasswordReset(user.ID)
	err = usc.passwordResetDB.Create(ctx, passwordRest)
	if err != nil {
		return err
	}

	template := `
		<h1>%s</h1>
		<br/>
		<a href="%s">Reset Password</a>
	`

	changePasswordUrl := strings.ReplaceAll(
		usc.config.UrlPageChangePassword,
		"{token}",
		passwordRest.Token,
	)
	subject := usc.config.Title + " - Reset Password"
	content := fmt.Sprintf(template, usc.config.Title, changePasswordUrl)

	usc.mailService.Send(mailTo, subject, content)
	if err != nil {
		return err
	}

	return nil
}

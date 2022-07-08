package recover

import (
	"context"

	"github.com/edmarfelipe/next-u/services/identity/infra"
	"github.com/edmarfelipe/next-u/services/identity/infra/db"
)

type Usecaser interface {
	Execute(ctx context.Context, input Input) error
}

type Usecase struct {
	userRepository db.UserRepositorier
	mailService    infra.MailServicer
	validator      infra.Validatorer
}

func NewUsecase(userRepository db.UserRepositorier, mailService infra.MailServicer, validator infra.Validatorer) Usecaser {
	return &Usecase{
		userRepository: userRepository,
		mailService:    mailService,
		validator:      validator,
	}
}

type Input struct {
	Email string `json:"email" validate:"email"`
}

func (usc Usecase) Execute(ctx context.Context, input Input) error {
	err := usc.validator.IsValid(input)
	if err != nil {
		return err
	}

	user, err := usc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	mailTo := infra.MailTo{
		Name:  user.Name,
		Email: user.Email,
	}

	content := `
		<h1>Recover Password</h1>
	`

	usc.mailService.Send(mailTo, "Recover Password", content)
	if err != nil {
		return err
	}

	// https://github.com/sendgrid/sendgrid-go#with-mail-helper-class
	// https://medium.com/@dhanushgopinath/sending-html-emails-using-templates-in-golang-9e953ca32f3d
	// https://github.com/sendgrid/email-templates/tree/master/paste-templates
	return nil
}

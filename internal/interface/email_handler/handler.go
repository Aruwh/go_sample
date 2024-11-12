package emailhandler

import (
	"fewoserv/internal/infrastructure/email_template"
	"net/smtp"
)

type (
	IEmailHandler interface {
		Send(toEmail string, emailTemplateHandler email_template.IEmailTemplate) error
	}

	SenderCredential struct {
		Email          string
		Password       string
		SmtpServerUrl  string
		SmtpServerPort string
	}

	EmailHandler struct {
		credentials SenderCredential
	}
)

func New(fromEmail, password, smtpServerUrl, smtpServerPort string) IEmailHandler {
	credentials := SenderCredential{
		Email:          fromEmail,
		Password:       password,
		SmtpServerUrl:  smtpServerUrl,
		SmtpServerPort: smtpServerPort,
	}

	emailHandler := EmailHandler{credentials: credentials}
	return &emailHandler
}

func (eh *EmailHandler) Send(toEmail string, emailTemplate email_template.IEmailTemplate) error {
	auth := smtp.PlainAuth("", eh.credentials.Email, eh.credentials.Password, eh.credentials.SmtpServerUrl)
	smtpServerAddress := eh.credentials.SmtpServerUrl + ":" + eh.credentials.SmtpServerPort

	emailMessage, err := emailTemplate.Build("FeWoServ", eh.credentials.Email, toEmail)
	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpServerAddress, auth, eh.credentials.Email, []string{toEmail}, emailMessage)
	if err != nil {
		return err
	}

	return nil
}

package mail

import (
	"errors"

	gomail "gopkg.in/mail.v2"
)

type genericMail struct {
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}

func NewGenericMail(
	host string,
	port int,
	user string,
	pass string,
) (Service, error) {
	if len(host) == 0 {
		return nil, errors.New("host is required")
	}
	if port <= 0 || port > 65535 {
		return nil, errors.New("port is required")
	}

	if len(user) == 0 {
		return nil, errors.New("user is required")
	}
	if len(pass) == 0 {
		return nil, errors.New("password is required")
	}
	return &genericMail{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
	}, nil
}

func (mail *genericMail) SendMail(from, to, subject, content string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", content)

	dialer := gomail.NewDialer(
		mail.SMTPHost,
		mail.SMTPPort,
		mail.SMTPUser,
		mail.SMTPPass,
	)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

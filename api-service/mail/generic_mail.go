package mail

import gomail "gopkg.in/mail.v2"

type GenericMail struct {
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
) *GenericMail {
	return &GenericMail{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
	}
}

func (mail *GenericMail) SendMail(from, to, subject, content string) error {
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

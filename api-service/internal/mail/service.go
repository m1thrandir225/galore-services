// Package mail handles all the required logic for sending e-mails to the users.
package mail

// Service defines a way to send mail
type Service interface {
	SendMail(from, to, subject, content string) error
}

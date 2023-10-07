package mailer

import (
	"crypto/tls"
	"errors"
	"fmt"
	"gopkg.in/mail.v2"
	"milky-mailer/internal/configer"
)

func Send(s configer.EmailSenderConfig, To, Subject, ContentType, Body string) error {

	email := mail.NewMessage()

	email.SetHeader("From", fmt.Sprintf("%s <%s>", s.FromName, s.From))
	email.SetHeader("To", To)
	email.SetHeader("Subject", Subject)

	email.SetBody(ContentType, Body)

	dialer := mail.NewDialer(s.Host, s.Port, s.User, s.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: !s.Tls, ServerName: s.Host}

	err := dialer.DialAndSend(email)
	if err != nil {
		return errors.Join(err, errors.New("SMTP: error send email"))
	}

	return nil
}

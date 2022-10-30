package mailer

import (
	"crypto/tls"
	"gopkg.in/mail.v2"
	"milky-mailer/internal/configer"
)

type EmailData struct {
	To          string
	FromName    string
	Subject     string
	ContentType string
	Body        string
}

func SendEmail(cfg *configer.EmailSenderConfig, data *EmailData) error {

	// TODO Валидация структур приходящих данных

	m := mail.NewMessage()

	m.SetHeader("From", data.FromName+" <"+cfg.From+">")
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)

	m.SetBody(data.ContentType, data.Body)

	d := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: cfg.Host}

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	return nil
}

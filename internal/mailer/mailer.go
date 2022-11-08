package mailer

import (
	"crypto/tls"
	"gopkg.in/mail.v2"
)

type EmailData struct {
	To          string
	FromName    string
	Subject     string
	ContentType string
	Body        string
}

type EmailSenderConfig struct {
	From     string
	User     string
	Host     string
	Password string
	Port     int
}

func SendEmail(cfg *EmailSenderConfig, data *EmailData) error {

	// TODO Валидация структур приходящих данных

	email := mail.NewMessage()

	email.SetHeader("From", data.FromName+" <"+cfg.From+">")
	email.SetHeader("To", data.To)
	email.SetHeader("Subject", data.Subject)

	email.SetBody(data.ContentType, data.Body)

	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)

	// TODO Надо сделать проверку на то, что в конфиге указано использовать SSL
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false, ServerName: cfg.Host}

	err := dialer.DialAndSend(email)
	if err != nil {
		return err
	}

	return nil
}

package mailer

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(host string, port int, username, password string) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &Mailer{dialer: dialer}
}

func (m *Mailer) SendMail(to, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", m.dialer.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return m.dialer.DialAndSend(message)
}
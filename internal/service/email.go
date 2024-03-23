package service

import (
	"crypto/tls"
	"log/slog"

	gomail "gopkg.in/mail.v2"
)

type EmailService struct {
	from     string
	password string
	smtpHost string
	smtpPort string
	log      *slog.Logger
}

func NewEmailService(from, password, smtpHost, smtpPort string, log *slog.Logger) *EmailService {

	return &EmailService{from, password, smtpHost, smtpPort, log}
}

func (s *EmailService) SendTextEmail(to, body, subject string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", s.from)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body)

	// Settings for SMTP server
	d := gomail.NewDialer(s.smtpHost, 587, s.from, s.password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		s.log.Debug(err.Error())
	}
}

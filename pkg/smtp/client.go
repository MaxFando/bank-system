package smtp

import (
	"crypto/tls"
	"fmt"
	"github.com/go-mail/mail/v2"
	"log"
)

const (
	smtpHost = "smtp.example.com"    // Хост SMTP-сервера
	smtpPort = 587                   // Порт (чаще используется 587 с TLS)
	smtpUser = "noreply@example.com" // Учетная запись
	smtpPass = "strong_password"     // Пароль/токен
)

func createMessage(to string, subject string, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return m
}

func createDialer() *mail.Dialer {
	d := mail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: false, // Не отключать проверку сертификата
	}
	return d
}

func sendEmail(d *mail.Dialer, m *mail.Message) error {
	if err := d.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("email sending failed")
	}
	return nil
}

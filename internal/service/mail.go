package service

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/spf13/viper"
)

const (
	activationLink = "http://%s:%s/auth/activate/%s"
	resetPasswordLink = "http://%s:%s/reset/%s"
)

type MailService struct {
	from string
	pass string
	host string
	port string
}

func NewMailService() *MailService {
	return &MailService{
		from: os.Getenv("EMAIL"),
		pass: os.Getenv("EMAIL_PASS"),
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
	}
}

func (s *MailService) SendActivationMail(to []string, link string) error {
	subject := "Account activation"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title></title>
		</head>
		<body>
			<h1>To activate your account, click the button below</h1>
			<a href="%s" style="padding: 12px 80px;background: #ffe057;color: #121212;text-decoration: none;border-radius: 50px;text-transform: uppercase;font-family: monospace;font-size: 18px;font-weight: 600;">Activate</a>
		</body>
		</html>
	`, fmt.Sprintf(activationLink, viper.GetString("app.host"), viper.GetString("app.port"), link))

	msg := []byte("Subject: " + subject + "\r\n" +
	"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
	"\r\n" + body)

	auth := smtp.PlainAuth("", s.from, s.pass, s.host)

	if err := smtp.SendMail(s.host + ":" + s.port, auth, s.from, to, msg); err != nil {
		return err
	}

	return nil
}

func (s *MailService) SendResetPasswordToken(to []string, link string) error {
	subject := "Password Reset"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title></title>
		</head>
		<body>
			<h1>To reset your password, click the button below</h1>
			<a href="%s" style="padding: 12px 80px;background: #ffe057;color: #121212;text-decoration: none;border-radius: 50px;text-transform: uppercase;font-family: monospace;font-size: 18px;font-weight: 600;">Reset</a>
		</body>
		</html>
	`, fmt.Sprintf(resetPasswordLink, viper.GetString("app.host"), viper.GetString("app.port"), link))

	msg := []byte("Subject: " + subject + "\r\n" +
	"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
	"\r\n" + body)

	auth := smtp.PlainAuth("", s.from, s.pass, s.host)

	if err := smtp.SendMail(s.host + ":" + s.port, auth, s.from, to, msg); err != nil {
		return err
	}

	return nil
}

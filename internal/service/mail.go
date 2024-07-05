package service

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/morf1lo/blog/internal/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	activationLink = "http://%s:%s/auth/activate/%s"
	resetPasswordLink = "http://%s:%s/reset/%s"
)

type MailService struct {
	rabbitMQ *mq.MQConn

	from string
	pass string
	host string
	port string
}

func NewMailService(rabbitMQ *mq.MQConn) *MailService {
	return &MailService{
		rabbitMQ: rabbitMQ,
		from: os.Getenv("EMAIL"),
		pass: os.Getenv("EMAIL_PASS"),
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
	}
}

func (s *MailService) ProcessActivationMails() {
	msgs, err := s.rabbitMQ.Consume("notifications.mail.activation")
	if err != nil {
		logrus.Fatalf("error starting consumer: %s", err.Error())
	}

	go func ()  {
		for msg := range msgs {
			var message ActivationMailData
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				logrus.Errorf("failed unmarshal message: %s", err.Error())
				msg.Nack(false, true)
				continue
			}

			if err := s.SendActivationMail(message.To, message.Link); err != nil {
				logrus.Errorf("failed send activation mail: %s", err.Error())
				msg.Nack(false, true)
				continue
			}

			if err := msg.Ack(false); err != nil {
				logrus.Errorf("failed ack message: %s", err.Error())
			}

			logrus.Println("Activation mail was sent successfully!")
			time.Sleep(time.Second * 3)
		}
	}()
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

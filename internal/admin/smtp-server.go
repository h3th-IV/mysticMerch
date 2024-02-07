package admin

import (
	"os"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"gopkg.in/gomail.v2"
)

func SendEmailNotification(recipient, subject, body string) error {
	return nil
}

type SMTPsErver struct {
	Host     string
	Port     int
	Username string
	Passowrd string
}

func NewSMTP() *SMTPsErver {
	return &SMTPsErver{
		Host:     "smtp.protonmail.com",
		Port:     465,
		Username: os.Getenv("NIMDALIAME"),
		Passowrd: os.Getenv("NIMDASSAP"),
	}
}
func TransactionalEmail(user *models.User, subject, body string) error {
	smtp := NewSMTP()
	Dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Passowrd)
	Mailer := gomail.NewMessage()
	Mailer.SetHeader("From", smtp.Username)
	Mailer.SetHeader("To", *user.Email)
	Mailer.SetBody("text/html", body)
	if err := Dialer.DialAndSend(Mailer); err != nil {
		return err
	}
	return nil
}

func MarketingEmail() {}

package admin

import (
	"os"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"gopkg.in/gomail.v2"
)

func SendEmailNotification(recipient, subject, body string) error {
	return nil
}

type SMTPServer struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTP() *SMTPServer {
	return &SMTPServer{
		Host:     "smtp.protonmail.com",
		Port:     465,
		Username: os.Getenv("NIMDALIAME"),
		Password: os.Getenv("NIMDASSAP"),
	}
}

// trabsactional email sent to each user concerning the state of their transaction
func TransactionalEmail(user *models.User, subject, body string) error {
	smtp := NewSMTP()
	Dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	Mailer := gomail.NewMessage()
	Mailer.SetHeader("From", smtp.Username)
	Mailer.SetHeader("To", *user.Email)
	Mailer.SetBody("text/html", body)
	if err := Dialer.DialAndSend(Mailer); err != nil {
		return err
	}
	return nil
}

// some form of Broadcast email
func MarketingEmail(users []*models.ResponseUser, subject, body string) error {
	smtp := NewSMTP()
	Dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	Mailer := gomail.NewMessage()
	Mailer.SetHeader("From", smtp.Username)
	for _, user := range users {
		Mailer.SetHeader("To", user.Email)
		Mailer.SetBody("text/html", body)
		if err := Dialer.DialAndSend(); err != nil {
			return err
		}
	}
	return nil
}

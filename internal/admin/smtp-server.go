package admin

import (
	"os"

	"github.com/h3th-IV/mysticMerch/internal/models"
	"github.com/h3th-IV/mysticMerch/internal/utils"
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
	if err := utils.LoadEnv(); err != nil {
		return nil
	}
	return &SMTPServer{
		Host:     "smtp.protonmail.com",
		Port:     465,
		Username: os.Getenv("NIMDALIAME"),
		Password: os.Getenv("NIMDASSAP"),
	}
}

// trabsactional email sent to each user concerning the state of their transaction
func TransactionalEmail(user *models.ResponseUser, subject, body string) error {
	smtp := NewSMTP()
	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtp.Username)
	mailer.SetHeader("To", user.Email)
	mailer.SetBody("text/html", body)
	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}

// some form of Broadcast email
func MarketingEmail(users []*models.ResponseUser, subject, body string) error {
	smtp := NewSMTP()
	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Username, smtp.Password)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtp.Username)
	for _, user := range users {
		mailer.SetHeader("To", user.Email)
		mailer.SetBody("text/html", body)
		if err := dialer.DialAndSend(); err != nil {
			return err
		}
	}
	return nil
}

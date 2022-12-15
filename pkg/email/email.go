package email

import (
	"fmt"
	"murakali/config"
	"net/smtp"
)

type Email struct {
	Auth smtp.Auth
	Addr string
	From string
}

func SendEmail(cfg *config.Config, toEmail, subject, msg string) {
	email := &Email{
		Auth: smtp.PlainAuth("", cfg.External.SMTPFrom, cfg.External.SMTPPassword, cfg.External.SMTPHost),
		Addr: fmt.Sprintf("%s:%v", cfg.External.SMTPHost, cfg.External.SMTPPort),
		From: cfg.External.SMTPFrom,
	}

	from := email.From
	value := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"mime := MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" + msg + "\r\n")

	err := smtp.SendMail(email.Addr, email.Auth, from, []string{toEmail}, value)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
}

package utils

import (
	"os"
	"gopkg.in/gomail.v2"
)


func SendConfirmationEmail(to string, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Confirm your email")
	m.SetBody("text/html", "Click <a href='http://localhost:8080/confirm/"+token+"'>here</a> to confirm your email")

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	return d.DialAndSend(m)
}




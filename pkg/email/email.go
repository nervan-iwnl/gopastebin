package email

import (
	"fmt"
	"net/smtp"

	"gopastebin/config"
)

func SendVerifyEmail(cfg *config.Config, to, code string) error {
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" || cfg.SMTPPass == "" {
		return nil // просто не шлём, но и не падаем
	}

	link := fmt.Sprintf("http://localhost:%s/api/v1/auth/verify?email=%s&code=%s", cfg.AppPort, to, code)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Verify your email\r\n" +
		"\r\n" +
		"Click here to verify: " + link + "\r\n")

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	return smtp.SendMail(cfg.SMTPHost+":"+cfg.SMTPPort, auth, cfg.SMTPUser, []string{to}, msg)
}

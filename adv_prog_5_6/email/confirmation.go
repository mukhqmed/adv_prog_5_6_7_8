package email

import (
	"fmt"
	"net/smtp"
)

func SendConfirmationEmail(email, token string) error {
	// Compose confirmation URL with token
	confirmationURL := fmt.Sprintf("http://localhost:8080/confirm?token=%s", token)

	// Compose email message
	message := fmt.Sprintf("Click the following link to confirm your registration: %s", confirmationURL)

	// Use an SMTP service to send the email
	// Update SMTP settings with your own
	smtpHost := "smtp.example.com"
	smtpPort := 587
	smtpUsername := "your_smtp_username"
	smtpPassword := "your_smtp_password"

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth,
		"from@example.com",
		[]string{email},
		[]byte(message),
	)
	if err != nil {
		return err
	}

	return nil
}

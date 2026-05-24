package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendOTPEmail uses standard SMTP to send a 6-digit code to the user
func SendOTPEmail(toEmail string, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	address := fmt.Sprintf("%s:%s", host, port)
	
	subject := "Subject: Your LendoGo Verification Code\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<h2>Welcome to LendoGo!</h2>
		<p>Your verification code is: <strong>%s</strong></p>
		<p>This code will expire in 5 minutes.</p>
	`, otp)

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, []string{toEmail}, message)
	
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
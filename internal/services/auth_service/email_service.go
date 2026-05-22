package auth

import (
	"fmt"
	"os" // <-- We need this to read your .env file!
	"gopkg.in/gomail.v2"
)

func SendOTPEmail(toEmail string, otpCode string) error {
	// 1. Grab your REAL credentials from the .env file
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	fmt.Printf("\n 	Debug:Email is '%s'\n",senderEmail)
	fmt.Printf("\n Debug password is '%s'(Lendth '%d')\n\n",senderPassword,len(senderPassword	))

	// 2. Configure the Email Message
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail) // Use the variable, not a string!
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your LendoGo Verification Code")
	
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; text-align: center; padding: 20px;">
			<h2>Welcome to LendoGo! 🚀</h2>
			<p>Your secure verification code is:</p>
			<h1 style="color: #0d6efd; letter-spacing: 5px;">%s</h1>
			<p>This code will expire in 5 minutes. Do not share this with anyone.</p>
		</div>
	`, otpCode)
	
	m.SetBody("text/html", htmlBody)

	// 3. Configure the SMTP Server using the REAL variables!
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPassword)

	// 4. Send it!
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
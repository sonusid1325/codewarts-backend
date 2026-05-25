package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// SendVerificationEmail sends a 6-digit verification code to the target email.
func SendVerificationEmail(targetEmail, username, code string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpSender := os.Getenv("SMTP_SENDER")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || smtpSender == "" {
		return fmt.Errorf("SMTP configuration incomplete: ensure SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS, SMTP_SENDER are set in environment")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Compose Email
	subject := "Subject: [Codewarts] Verify Your Wand Coordinates!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: sans-serif; background-color: #09090b; color: #f4f4f5; padding: 20px; margin: 0;">
			<div style="max-width: 600px; margin: 0 auto; background-color: #18181b; border: 1px solid #27272a; padding: 40px; border-radius: 16px; box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);">
				<div style="text-align: center; margin-bottom: 30px;">
					<h1 style="color: #fbbf24; margin: 0; font-size: 30px; letter-spacing: 2px; font-weight: 800;">CODEWARTS</h1>
					<p style="color: #a1a1aa; font-size: 13px; margin: 5px 0 0 0; text-transform: uppercase; letter-spacing: 1.5px; font-family: monospace;">Chamber of Shells Protocol</p>
				</div>
				<div style="font-size: 16px; line-height: 1.7; color: #e4e4e7;">
					<p>Hail, Wizard <strong>%s</strong>,</p>
					<p>Snape has locked the dungeon mainframe. To align your wand coordinates and access your private sandbox container, you must verify your identity with the following verification rune:</p>
					<div style="text-align: center; margin: 35px 0;">
						<span style="display: inline-block; background-color: #09090b; border: 2px dashed #fbbf24; color: #fbbf24; font-size: 36px; font-family: monospace; font-weight: bold; letter-spacing: 8px; padding: 18px 36px; border-radius: 12px; box-shadow: 0 0 20px rgba(251, 191, 36, 0.05);">
							%s
						</span>
					</div>
					<p style="font-size: 14px; color: #71717a; line-height: 1.6; border-left: 3px solid #d946ef; padding-left: 12px; margin-top: 25px;">
						This code will expire shortly. Do not share your wand passkey or verification runes with other students, including members of Slytherin.
					</p>
				</div>
				<hr style="border: 0; border-top: 1px solid #27272a; margin: 35px 0;">
				<div style="text-align: center; font-size: 11px; color: #52525b; font-family: monospace; letter-spacing: 1px;">
					CODEWARTS // WIZARDING SHELL ACADEMY
				</div>
			</div>
		</body>
		</html>
	`, username, code)

	msg := []byte("To: " + targetEmail + "\n" + "From: Codewarts <" + smtpSender + ">\n" + subject + mime + body)
	addr := smtpHost + ":" + smtpPort

	log.Printf("Connecting to SMTP server at %s for %s...", addr, targetEmail)
	c, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer c.Close()

	config := &tls.Config{
		ServerName: smtpHost,
	}
	if err = c.StartTLS(config); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate SMTP: %w", err)
	}

	if err = c.Mail(smtpSender); err != nil {
		return fmt.Errorf("failed to execute MAIL command: %w", err)
	}

	if err = c.Rcpt(targetEmail); err != nil {
		return fmt.Errorf("failed to execute RCPT command: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to execute DATA command: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	if err = c.Quit(); err != nil {
		log.Printf("Warning: SMTP QUIT connection error: %v", err)
	}

	log.Printf("Verification email successfully delivered to %s", targetEmail)
	return nil
}

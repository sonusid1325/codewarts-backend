package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// SendVerificationEmail sends a 6-digit verification code to the target email.
func SendVerificationEmail(targetEmail, username, code string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	smtpSender := os.Getenv("SMTP_SENDER")

	if apiKey == "" || smtpSender == "" {
		return fmt.Errorf("Brevo configuration incomplete: ensure BREVO_API_KEY and SMTP_SENDER are set in environment")
	}

	// Compose Email
	subject := "[Codewarts] Verify Your Wand Coordinates!"
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

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Codewarts",
			"email": smtpSender,
		},
		"to": []map[string]string{
			{
				"email": targetEmail,
				"name":  username,
			},
		},
		"subject":     subject,
		"htmlContent": body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode Brevo payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.brevo.com/v3/smtp/email", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create Brevo request: %w", err)
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Dispatching verification email to %s via Brevo API...", targetEmail)
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Brevo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo send failed: status %s: %s", resp.Status, strings.TrimSpace(string(bodyBytes)))
	}

	log.Printf("Verification email successfully delivered to %s", targetEmail)
	return nil
}

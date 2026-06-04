package mailer

import (
	"context"
	"log"
	"os"

	"github.com/resend/resend-go/v3"
)

var client *resend.Client
var sender string

func Init() {
	resend_api_key := os.Getenv("RESEND_API_KEY")
	if resend_api_key == "" {
		log.Println("WARNING: Resend API key not set.")
	}
	sender := os.Getenv("SENDER_EMAIL")
	if sender == "" {
		log.Println("WARNING: Sender email not set.")
	}
	client = resend.NewClient(resend_api_key)
}

func SendEmailAsync(to []string, subject string, htmlContent string) {
	go func() {
		params := &resend.SendEmailRequest{
			From:    sender,
			To:      to,
			Html:    htmlContent,
			Subject: subject,
		}

		_, err := client.Emails.SendWithContext(context.Background(), params)
		if err != nil {
			log.Printf("Failed to send email to %v: %v", to, err)
			return
		}

		log.Printf("Email successfully sent asynchronously to %v", to)
	}()
}

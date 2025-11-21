package service

import (
	"fmt"
	"os"
)

// WhatsAppService handles sending WhatsApp messages via Twilio API
type WhatsAppService struct {
	AccountSID string
	AuthToken  string
	FromNumber string // Twilio WhatsApp number (format: whatsapp:+14155238886)
}

// NewWhatsAppService creates a new WhatsApp service instance
func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		FromNumber: os.Getenv("TWILIO_WHATSAPP_FROM"),
	}
}

// SendMessage sends a WhatsApp message to a phone number
func (w *WhatsAppService) SendMessage(to, message string) error {
	// Check if WhatsApp is configured
	if w.AccountSID == "" || w.AuthToken == "" || w.FromNumber == "" {
		return fmt.Errorf("WhatsApp service not configured. Set TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_WHATSAPP_FROM env vars")
	}

	// Format the 'to' number (ensure it's in whatsapp: format)
	toNumber := to
	if len(to) > 0 && to[0] != '+' {
		toNumber = "+" + to
	}
	if len(toNumber) < 10 {
		toNumber = "whatsapp:" + toNumber
	} else if toNumber[:10] != "whatsapp:+" {
		toNumber = "whatsapp:" + toNumber
	}

	// For MVP, we'll log the message (actual Twilio API call would go here)
	// TODO: Implement actual Twilio API call when API keys are available
	fmt.Printf("[WhatsApp] Would send to %s: %s\n", toNumber, message)

	// Example Twilio API call (commented out - requires http client):
	// url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", w.AccountSID)
	// data := url.Values{}
	// data.Set("From", w.FromNumber)
	// data.Set("To", toNumber)
	// data.Set("Body", message)
	//
	// req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	// req.SetBasicAuth(w.AccountSID, w.AuthToken)
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	//     return err
	// }
	// defer resp.Body.Close()

	return nil
}

// SendEmergencyAlert sends emergency alerts to multiple recipients
func (w *WhatsAppService) SendEmergencyAlert(recipients []string, message string) map[string]error {
	results := make(map[string]error)
	for _, recipient := range recipients {
		if recipient == "" {
			continue
		}
		err := w.SendMessage(recipient, message)
		results[recipient] = err
	}
	return results
}

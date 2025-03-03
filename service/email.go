package service

import (
	"fmt"
	"net/smtp"
	"os"
)

// EmailConfig menyimpan konfigurasi untuk pengiriman email.
type EmailConfig struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	Recipient    string
}

// LoadEmailConfig memuat konfigurasi email dari variabel lingkungan.
func LoadEmailConfig() (*EmailConfig, error) {
	server := os.Getenv("SMTP_SERVER")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	recipient := os.Getenv("ALERT_EMAIL")
	if server == "" || port == "" || username == "" || password == "" || recipient == "" {
		return nil, fmt.Errorf("konfigurasi email tidak lengkap, periksa SMTP_SERVER, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, ALERT_EMAIL")
	}
	return &EmailConfig{
		SMTPServer:   server,
		SMTPPort:     port,
		SMTPUsername: username,
		SMTPPassword: password,
		Recipient:    recipient,
	}, nil
}

// SendNotification mengirimkan notifikasi email ketika mode berubah.
func (config *EmailConfig) SendNotification(user, zoneID, mode string, cpuUsage float64) error {
	var subject string
	if mode == "on" {
		subject = "Under Attack Mode Activated"
	} else {
		subject = "Under Attack Mode Deactivated"
	}

	body := fmt.Sprintf(
		"Notifikasi perubahan mode:\n\nUser: %s\nZone: %s\nCPU Usage: %.2f%%\nAction: %s",
		user, zoneID, cpuUsage, subject,
	)

	// Susun pesan email.
	msg := "From: " + config.SMTPUsername + "\n" +
		"To: " + config.Recipient + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPServer)
	addr := fmt.Sprintf("%s:%s", config.SMTPServer, config.SMTPPort)

	// Kirim email.
	err := smtp.SendMail(addr, auth, config.SMTPUsername, []string{config.Recipient}, []byte(msg))
	if err != nil {
		return fmt.Errorf("gagal mengirim email: %v", err)
	}
	return nil
}

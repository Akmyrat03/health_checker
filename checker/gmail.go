package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type SMTPConfig struct {
	SMTPServer    string   `json:"smtp_server"`
	SMTPPort      string   `json:"smtp_port"`
	SMTPEmail     string   `json:"smtp_email"`
	SMTPPass      string   `json:"smtp_pass"`
	SubjectPrefix string   `json:"subject_prefix"`
	Users         []string `json:"users"`
}

type Config struct {
	SMTPConfig SMTPConfig `json:"smtp_config"`
}

// LoadEmailConfig loads the email configuration from a JSON file
func LoadEmailConfig(filename string) (*SMTPConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open email config file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read email config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse email config JSON: %v", err)
	}

	return &config.SMTPConfig, nil
}

type EmailSender struct {
	Config *SMTPConfig
}

func NewEmailSender(configFile string) (*EmailSender, error) {
	config, err := LoadEmailConfig(configFile)
	if err != nil {
		return nil, err
	}
	return &EmailSender{Config: config}, nil
}

func (s *EmailSender) SendError(message string) error {
	auth := smtp.PlainAuth("", s.Config.SMTPEmail, s.Config.SMTPPass, s.Config.SMTPServer)

	usersTo := s.Config.Users
	subject := s.Config.SubjectPrefix + " - Error Notification"
	body := message

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\n", s.Config.SMTPEmail))
	msg.WriteString(fmt.Sprintf("To: %s\n", strings.Join(usersTo, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\n\n", subject))
	msg.WriteString(body)

	emailMessage := msg.String()

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", s.Config.SMTPServer, s.Config.SMTPPort),
		auth,
		s.Config.SMTPEmail,
		usersTo,
		[]byte(emailMessage),
	)

	if err != nil {
		return fmt.Errorf("failed to send email notification: %w", err)
	}

	return nil
}

func SendEmailNotification(message string) {
	emailSender, err := NewEmailSender("config.json")
	if err != nil {
		log.Printf("Failed to load email config: %v", err)
		return
	}

	err = emailSender.SendError(message)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	}
}

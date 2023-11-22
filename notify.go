package main

import (
	"log"
	"net/smtp"
)

var (
	AdminEmail    = "admin@example.com"
	EmailServer   = "smtp.gmail.com:587"
	EmailUser     = "notify@example.com"
	EmailPassword = "password"
)

// NotifyAdmin sends an email notification to the admin about the failure of an event
func NotifyAdmin(subject, body string) {
	from := EmailUser
	to := []string{AdminEmail}
	msg := "From: " + from + "\n" +
		"To: " + AdminEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(EmailServer, smtp.PlainAuth("", EmailUser, EmailPassword, "smtp.gmail.com"), from, to, []byte(msg))
	if err != nil {
		log.Printf("Error notifying the admin: %v", err)
	}

}

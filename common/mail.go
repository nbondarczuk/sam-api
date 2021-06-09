package common

import (
	"bytes"
	"log"
	"net/smtp"
)

func MailTo(address, message string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(AppConfig.AlertMailServerAddress)
	if err != nil {
		log.Printf("Error dialing mail server: %s %s", AppConfig.AlertMailServerAddress, err.Error())
		return err
	}
	defer c.Close()

	// Set the sender and recipient.
	c.Mail(AppConfig.AlertMailSenderAddress)
	c.Rcpt(address)

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Printf("Error sending mail body to server: %s %s", AppConfig.AlertMailServerAddress, err.Error())
		return err
	}
	defer wc.Close()
	buf := bytes.NewBufferString(message)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Printf("Error writing: %s %s", AppConfig.AlertMailServerAddress, err.Error())
		return err
	}

	return nil
}

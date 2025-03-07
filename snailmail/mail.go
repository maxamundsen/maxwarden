// This package wraps around the stdlib's smtp client implementation, allowing
// for easily sending emails. This package should be used alongside the "views"
// package, to generate an HTML template for use in an email.
package snailmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"net/smtp"
	"maxwarden/config"
	"strings"
	"time"
)

const (
	TYPE_TEXT = iota
	TYPE_HTML = iota
)

type Email struct {
	Recipients []string
	Subject    string
	Body       *bytes.Buffer
}

func SendMail(message Email, mailtype int) error {
	recipientString := strings.Join(message.Recipients, ",")
	from := mail.Address{Name: config.GetConfig().SmtpDisplayFrom, Address: config.GetConfig().SmtpUsername}

	header := make(map[string]string)
	header["To"] = recipientString
	header["From"] = from.String()
	header["Subject"] = mime.QEncoding.Encode("UTF-8", message.Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Transfer-Encoding"] = "base64"
	header["Date"] = time.Now().Format(time.RFC1123)

	if mailtype == TYPE_HTML {
		header["Content-Type"] = "text/html; charset=\"utf-8\""
	} else {
		header["Content-Type"] = "text/plain; charset=\"utf-8\""
	}

	email := ""
	for k, v := range header {
		email += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	email += "\r\n" + base64.StdEncoding.EncodeToString(message.Body.Bytes())

	var auth smtp.Auth = nil

	if config.GetConfig().SmtpRequireAuth {
		auth = smtp.PlainAuth("", config.GetConfig().SmtpUsername, config.GetConfig().SmtpPassword, config.GetConfig().SmtpServer)
	}

	err := smtp.SendMail(config.GetConfig().SmtpServer+":"+config.GetConfig().SmtpPort, auth, config.GetConfig().SmtpUsername, message.Recipients, []byte(email))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

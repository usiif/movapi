package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dilar  *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dilar := mail.NewDialer(host, port, username, password)
	dilar.Timeout = 5 * time.Second

	return Mailer{
		dilar:  dilar,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetHeader("text/plain", plainBody.String())
	msg.SetHeader("text/html", htmlBody.String())

	for i := 1; i < 3; i++ {
		err = m.dilar.DialAndSend(msg)
		if nil != err {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

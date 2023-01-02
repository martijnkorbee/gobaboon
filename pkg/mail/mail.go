// Package mail contains a simple mailer that connects to your mail service and sends emails through the jobs channel or direct methods calls.
// At the moment only supports SMTP.
package mail

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
)

type Mailer struct {
	// settings holds all settings to connect and send email
	settings MailerSettings

	// Service is the type of mail service to use i.e. SMTP, API
	Service string

	// Templates is the full path to the email templates
	Templates string

	// Jobs is the jobs channel
	Jobs chan Message

	// Results is the results channel for the jobs channel
	Results chan Result
}

type MailerSettings struct {
	Domain     string
	Host       string
	Port       int
	Username   string
	Password   string
	AuthMethod string
	Encryption string
	From       string
	FromName   string
}

type Message struct {
	From        string
	FromName    string
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

type Result struct {
	Success bool
	Error   error
}

func NewMailer(settings MailerSettings, service string, path string) (*Mailer, error) {
	switch strings.ToLower(service) {
	case "smtp":
		return &Mailer{
			settings:  settings,
			Service:   service,
			Templates: path,
			Jobs:      make(chan Message, 20),
			Results:   make(chan Result, 20),
		}, nil
	default:
		return nil, errors.New("non supported mailer service")
	}
}

func (m *Mailer) ListenForMail() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mailer) Send(msg Message) error {
	switch m.Service {
	case "SMTP":
		return m.SendSMTPMessage(msg)
	default:
		// no mailer specified
		return errors.New("none or invalid mailer specified in .env file")
	}
}

func (m *Mailer) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	// load a html template
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// execute and store in a io buffer
	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.Data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func (m *Mailer) buildPlainTextMessage(msg Message) (string, error) {

	templateToRender := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)

	// load a template
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// execute and store in a io buffer
	var tpl bytes.Buffer
	err = t.ExecuteTemplate(&tpl, "body", msg.Data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

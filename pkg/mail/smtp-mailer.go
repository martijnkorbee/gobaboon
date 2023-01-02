package mail

import (
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func (m *Mailer) SendSMTPMessage(msg Message) error {
	formattedMsg, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainTextMsg, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.settings.Host
	server.Port = m.settings.Port
	server.Username = m.settings.Username
	server.Password = m.settings.Password
	// set encryption
	switch strings.ToLower(m.settings.Encryption) {
	case "tls":
		server.Encryption = mail.EncryptionSTARTTLS
	case "ssl":
		server.Encryption = mail.EncryptionSSLTLS
	case "none":
		server.Encryption = mail.EncryptionNone
	default:
		server.Encryption = mail.EncryptionSTARTTLS
	}
	// set authentication
	switch strings.ToLower(m.settings.AuthMethod) {
	case "login":
		server.Authentication = mail.AuthLogin
	case "CRAMMD5":
		server.Authentication = mail.AuthCRAMMD5
	default:
		server.Authentication = mail.AuthPlain
	}

	server.KeepAlive = false

	server.ConnectTimeout = 15 * time.Second
	server.SendTimeout = 15 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To...).
		SetSubject(msg.Subject).
		AddCc(msg.CC...).
		AddBcc(msg.BCC...).
		SetBody(mail.TextHTML, formattedMsg).
		AddAlternative(mail.TextPlain, plainTextMsg)

	// TODO: add attachments

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

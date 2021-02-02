package email

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/smtp"
)

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	subjPasswordRecovery = "Восстановление пароля" // nolint: gosec
	tplPasswordRecovery  = "passwordRecovery.html" // nolint: gosec
)

type Email struct {
	TplPath string
	Host    string
	Port    string
	From    string
	Pass    string
}

func New(tplPath, host, port, pass, from string) *Email {
	return &Email{
		TplPath: tplPath,
		Host:    host,
		Port:    port,
		From:    from,
		Pass:    pass,
	}
}

func (e *Email) SendRecoveryCodeEmail(to, username, code string) error {
	templateData := struct {
		Username     string
		RecoveryCode string
	}{
		Username:     html.UnescapeString(username),
		RecoveryCode: code,
	}

	tpl, err := e.parseTemplate(fmt.Sprintf("%s/%s", e.TplPath, tplPasswordRecovery), templateData)
	if err != nil {
		return err
	}

	body := "To: " + to + "\r\nSubject: " + subjPasswordRecovery + "\r\n" + mime + "\r\n" + tpl

	return e.send(to, body)
}

func (e *Email) send(to, body string) error {
	var (
		auth   = smtp.PlainAuth("", e.From, e.Pass, e.Host)
		server = fmt.Sprintf("%s:%s", e.Host, e.Port)
	)
	return smtp.SendMail(server, auth, e.From, []string{to}, []byte(body))
}

func (e *Email) parseTemplate(filename string, data interface{}) (string, error) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

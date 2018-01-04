package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const (
	Version = "0.0.1"
)

type MailService interface {
	SendText(subject string, body string, to ...string) error
	SendHtml(subject string, body string, to ...string) error
	SendTpl(subject string, tplfile string, data map[string]interface{}, to ...string) error
	UpdateConfig(Config)
}

type mailer struct {
	config        Config
	from          mail.Address
	auth          smtp.Auth
	authenticated bool
}

func New(cfg Config) MailService {
	m := &mailer{config: cfg}
	from := m.config.From
	if from == "" {
		from = m.config.Username
	}
	if cfg.FromAlias == "" {
		if cfg.Username != "" && strings.Contains(cfg.Username, "@") {
			m.from = mail.Address{Name: cfg.Username[0:strings.IndexByte(cfg.Username, '@')], Address: from}
		}
	} else {
		m.from = mail.Address{Name: cfg.FromAlias, Address: from}
	}
	return m
}

func (m *mailer) UpdateConfig(cfg Config) {
	m.config = cfg
}

func (m *mailer) SendText(subject string, body string, to ...string) error {
	return m.send(subject, body, false, to)
}

func (m *mailer) SendHtml(subject string, body string, to ...string) error {
	return m.send(subject, body, true, to)
}

func (m *mailer) SendTpl(subject string, tplfile string, data map[string]interface{}, to ...string) error {
	tplPath := m.config.TplPath
	if tplPath == "" {
		return fmt.Errorf("TplPath is empty ")
	}
	tpl, err := template.ParseFiles(filepath.Join(tplPath, tplfile))
	if err != nil {
		return err
	}
	strBuf := new(bytes.Buffer)
	err = tpl.Execute(strBuf, data)
	if err != nil {
		return err
	}
	return m.send(subject, strBuf.String(), true, to)
}

func (m *mailer) send(subject string, body string, isHtml bool, to []string) error {
	if !m.authenticated {
		cfg := m.config
		if !cfg.IsValid() {
			return fmt.Errorf("paramter cannot be empty !")
		}
		if cfg.AuthType == AUTH_TYPE_LOGIN {
			m.auth = LoginAuth(cfg.Username, cfg.Password, cfg.Host)
		} else if cfg.AuthType == AUTH_TYPE_CRAMMD5 {
			m.auth = smtp.CRAMMD5Auth(cfg.Username, cfg.Password)
		} else {
			m.auth = PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		}
		m.authenticated = true
	}
	host := fmt.Sprintf("%s:%s", m.config.Host, strconv.Itoa(m.config.Port))

	header := make(map[string]string)
	header["From"] = m.from.String()
	header["To"] = strings.Join(to, ",")

	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	if isHtml {
		header["Content-Type"] = "text/html; charset=UTF-8"
	} else {
		header["Content-Type"] = "text/plain; charset=UTF-8"
	}
	header["Content-Transfer-Encoding"] = "base64"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	return smtp.SendMail(host, m.auth, m.config.Username, to, []byte(message))
}

/**
若要解决超时问题，请参考
https://github.com/tryor/commons/tree/master/ismtp
*/

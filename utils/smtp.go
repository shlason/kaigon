package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"
	"text/template"

	"github.com/shlason/kaigon/configs"
)

func SendEmail(to []string, title string, templateFileName string, templatParams any) error {
	var body bytes.Buffer

	cwd, _ := os.Getwd()
	t, _ := template.ParseFiles(filepath.Join(cwd, fmt.Sprintf("./templates/%s", templateFileName)))
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", title, mimeHeaders)))
	err := t.Execute(&body, templatParams)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", configs.Smtp.Sender, configs.Smtp.Password, configs.Smtp.Host)
	err = smtp.SendMail(configs.Smtp.Host+":"+configs.Smtp.Port, auth, configs.Smtp.Sender, to, body.Bytes())

	return err
}

package mail

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/smtp"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
)

const (
	sendGridHost = "smtp.sendgrid.net"
	sendGridPort = "465"
	emailAddress = "contact@digicre.net"
	fromName     = "デジクリ"
)

func sendEmail(to string, subject string, body string) error {
	templateData := map[string]string{
		"address": to,
	}

	renderedBody, err := renderTemplate(body, templateData)
	if err != nil {
		return fmt.Errorf("テンプレートの展開に失敗しました: %w", err)
	}

	message := buildEmailMessage(to, subject, renderedBody)

	auth := smtp.PlainAuth("", emailAddress, env.SendGridApiKey, sendGridHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         sendGridHost,
	}

	conn, err := tls.Dial("tcp", sendGridHost+":"+sendGridPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("SMTP接続に失敗しました: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, sendGridHost)
	if err != nil {
		return fmt.Errorf("SMTPクライアントの作成に失敗しました: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP認証に失敗しました: %w", err)
	}

	if err := client.Mail(emailAddress); err != nil {
		return fmt.Errorf("送信元の設定に失敗しました: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("送信先の設定に失敗しました: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("データストリームの作成に失敗しました: %w", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("メール本文の書き込みに失敗しました: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("データストリームのクローズに失敗しました: %w", err)
	}

	logrus.Infof("メール送信成功: %s", to)
	return nil
}

func buildEmailMessage(to string, subject string, body string) string {
	var message strings.Builder

	encodedSubject := mime.QEncoding.Encode("UTF-8", subject)
	encodedFromName := mime.QEncoding.Encode("UTF-8", fromName)

	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", encodedFromName, emailAddress))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Reply-To: %s\r\n", emailAddress))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", encodedSubject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("\r\n")

	message.WriteString(body)

	return message.String()
}

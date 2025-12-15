package mail

import (
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

const (
	emailAddress = "contact@digicre.net"
	fromName     = "芝浦工業大学 デジクリ"
)

func sendEmail(to string, subject string, body string) error {
	if env.SendGridApiKey == "" {
		return fmt.Errorf("SENDGRID_API_KEYが設定されていません")
	}

	templateData := map[string]string{
		"address": to,
	}

	renderedBody, err := renderTemplate(body, templateData)
	if err != nil {
		return fmt.Errorf("テンプレートの展開に失敗しました: %w", err)
	}

	// SendGrid SDKを使用してメールを作成
	from := mail.NewEmail(fromName, emailAddress)
	toEmail := mail.NewEmail("", to)
	content := mail.NewContent("text/plain", renderedBody)
	message := mail.NewV3MailInit(from, subject, toEmail, content)

	// SendGridクライアントを作成して送信
	client := sendgrid.NewSendClient(env.SendGridApiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("メール送信に失敗しました: %w", err)
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		logrus.Infof("メール送信成功: %s (ステータスコード: %d)", to, response.StatusCode)
		return nil
	}

	return fmt.Errorf("メール送信に失敗しました (ステータスコード: %d): %s", response.StatusCode, response.Body)
}

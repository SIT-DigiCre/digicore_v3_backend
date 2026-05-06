package mail

import (
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/resend/resend-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	emailAddress = "contact@digicre.net"
	fromName     = "芝浦工業大学 デジクリ"
)

func SendEmail(to string, subject string, body string) error {
	if env.ResendApiKey == "" {
		return fmt.Errorf("RESEND_API_KEYが設定されていません")
	}

	templateData := map[string]string{
		"address": to,
	}

	renderedBody, err := renderTemplate(body, templateData)
	if err != nil {
		return fmt.Errorf("テンプレートの展開に失敗しました: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:fmt.Sprintf("%s <%s>", fromName, emailAddress),
		To:[]string{to},
		Subject:subject,
		Text:renderedBody,
	}

	// Resendクライアントを作成して送信
	client := resend.NewClient(env.ResendApiKey)
	response, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("メール送信に失敗しました: %w", err)
	}

	logrus.Infof("メール送信成功: %s (id: %s)", to, response.Id)
	return nil

}

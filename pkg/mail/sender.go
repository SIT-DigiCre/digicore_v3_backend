package mail

import (
	"context"
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/resend/resend-go/v3"
	"github.com/sirupsen/logrus"
)

const (
	emailAddress = "contact@digicre.net"
	fromName     = "芝浦工業大学 デジクリ"
)

func buildSendRequest(to string, subject string, body string) (*resend.SendEmailRequest, error) {
	templateData := map[string]string{
		"address": to,
	}

	renderedBody, err := renderTemplate(body, templateData)
	if err != nil {
		return nil, fmt.Errorf("テンプレートの展開に失敗しました: %w", err)
	}

	return &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", fromName, emailAddress),
		To:      []string{to},
		Subject: subject,
		Text:    renderedBody,
	}, nil
}

// 既存の実装を維持するための関数
func SendEmail(to string, subject string, body string) error {
	if env.ResendApiKey == "" {
		return fmt.Errorf("RESEND_API_KEYが設定されていません")
	}
	req, err := buildSendRequest(to, subject, body)
	if err != nil {
		return err
	}
	// Resendクライアントを作成して送信
	client := resend.NewClient(env.ResendApiKey)
	response, err := client.Emails.Send(req)
	if err != nil {
		return fmt.Errorf("メール送信に失敗しました: %w", err)
	}

	logrus.Infof("メール送信成功: %s (id: %s)", to, response.Id)
	return nil
}

func SendEmails(to []string, subject string, body string, ctx context.Context) (*resend.BatchEmailResponse, error) {
	if env.ResendApiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEYが設定されていません")
	}
	reqs := make([]*resend.SendEmailRequest, 0, len(to))
	for _, to := range to {
		req, err := buildSendRequest(to, subject, body)
		if err != nil {
			return nil, err
		}
		reqs = append(reqs, req)
	}
	client := resend.NewClient(env.ResendApiKey)
	return client.Batch.SendWithOptions(ctx, reqs, &resend.BatchSendEmailOptions{BatchValidation: resend.BatchValidationPermissive})
}

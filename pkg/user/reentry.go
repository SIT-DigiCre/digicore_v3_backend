package user

import (
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/mail"
	"github.com/sirupsen/logrus"
)

const maxReentryCount = 2

type reentryCount struct {
	ReentryCount int `db:"reentry_count"`
}

type reentryPendingCount struct {
	PendingCount int `db:"pending_count"`
}

type reentryDetail struct {
	ReentryId    string `db:"reentry_id"`
	UserId       string `db:"user_id"`
	ReentryCount int    `db:"reentry_count"`
	Status       string `db:"status"`
	Note         string `db:"note"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}

func notifyReentryApplied(studentNumber string) {
	address := fmt.Sprintf("%s@shibaura-it.ac.jp", studentNumber)
	if err := mail.SendEmail(address, "再入部申請受付のお知らせ", "再入部申請を受け付けました。運営側で確認後、結果をメールでお知らせします。"); err != nil {
		logrus.Warnf("再入部申請受付メール送信に失敗しました(%s): %v", address, err)
	}
}

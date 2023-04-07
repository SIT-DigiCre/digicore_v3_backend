package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
)

func GetSchoolYear() int {
	now := time.Now()
	month := int(now.Month())
	if 1 <= month && month <= 3 {
		return now.Year() - 1
	}
	return now.Year()
}

func GetYear() int {
	now := time.Now()
	return now.Year()
}

func GetAfterDate(year int, month int, day int) string {
	now := time.Now()
	now = now.AddDate(year, month, day)
	return now.Format("2006-01-02")
}

func GetUniqueString(str []string) []string {
	m := make(map[string]bool)
	uniq := []string{}

	for _, ele := range str {
		if !m[ele] {
			m[ele] = true
			uniq = append(uniq, ele)
		}
	}
	return uniq
}

func RenewalActiveLimit(dbClient db.TransactionClient, userId string, activeLimit string) *response.Error {
	params := struct {
		UserId      string `twowaysql:"userId"`
		ActiveLimit string `twowaysql:"activeLimit"`
	}{
		UserId:      userId,
		ActiveLimit: activeLimit,
	}
	_, err := dbClient.Exec("sql/utils/update_user_active_limit.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}

func NoticeMattermost(text string) {
	if env.MattermostWebHookURL == "" {
		logrus.Error("Not set mattermost web hook url")
		return
	}
	payload := struct {
		Text string `json:"text"`
	}{Text: text}
	p, err := json.Marshal(payload)
	if err != nil {
		logrus.Error(err)
		return
	}
	resp, err := http.Post(env.MattermostWebHookURL, "application/json", bytes.NewBuffer(p))
	if err != nil {
		logrus.Error(err)
		return
	}
	defer resp.Body.Close()
}

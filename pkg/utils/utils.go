package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
)

// 現在の年度を取得する。4月から翌年3月までを年度とする。
func GetSchoolYear() int {
	now := time.Now()
	month := int(now.Month())
	if 1 <= month && month <= 3 {
		return now.Year() - 1
	}
	return now.Year()
}

// 3月中の部費振込に対応するために、3月を翌年度として現在の年度を取得する。
func GetFiscalYear() int {
	now := time.Now()
	month := int(now.Month())
	if 1 <= month && month <= 2 {
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

func CalculateSchoolGradeFromStudentNumber(studentNumber string) (int, error) {
	if len(studentNumber) < 4 {
		return 0, fmt.Errorf("student number is too short: %s", studentNumber)
	}

	currentSchoolYear := GetSchoolYear()
	enterYear, err := strconv.Atoi(studentNumber[2:4])
	if err != nil {
		return 0, fmt.Errorf("student number has invalid enter year: %s", studentNumber)
	}

	schoolGrade := currentSchoolYear - 2000 - enterYear + 1
	switch studentNumber[0] {
	case 'm':
		schoolGrade += 4
	case 'n':
		schoolGrade += 6
	}

	return schoolGrade, nil
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

func NoticeMattermost(text string, channel string, username string, iconEmoji string) {
	if env.MattermostWebHookURL == "" {
		logrus.Error("Not set mattermost web hook url")
		return
	}
	payload := struct {
		Text      string `json:"text"`
		Channel   string `json:"channel"`
		Username  string `json:"username"`
		IconEmoji string `json:"icon_emoji"`
	}{Text: text, Channel: channel, Username: username, IconEmoji: iconEmoji}
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logrus.WithError(cerr).Warn("Mattermost Webhookレスポンスのクローズに失敗しました")
		}
	}()
}

type FileId struct {
	FileId string `db:"file_id"`
}

type FileInfo struct {
	FileId string `db:"file_id"`
	Name   string `db:"name"`
}

func GetFileInfo(dbClient db.Client, fileIds []string) (map[string]FileInfo, *response.Error) {
	if len(fileIds) == 0 {
		return map[string]FileInfo{}, nil
	}
	params := struct {
		FileIds []string `twowaysql:"fileIds"`
	}{
		FileIds: fileIds,
	}
	fileInfos := []FileInfo{}
	err := dbClient.Select(&fileInfos, "sql/utils/select_file.sql", &params)
	if err != nil {
		return map[string]FileInfo{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルの取得に失敗しました", Log: err.Error()}
	}
	res := map[string]FileInfo{}
	for _, v := range fileInfos {
		res[v.FileId] = v
	}
	return res, nil
}

type UserID struct {
	UserId string `db:"user_id"`
}

type UserInfo struct {
	IconUrl  string `db:"icon_url"`
	UserId   string `db:"user_id"`
	Username string `db:"username"`
}

func GetUserInfo(dbClient db.Client, userIds []string) (map[string]UserInfo, *response.Error) {
	if len(userIds) == 0 {
		return map[string]UserInfo{}, nil
	}
	params := struct {
		UserIds []string `twowaysql:"userIds"`
	}{
		UserIds: userIds,
	}
	userInfos := []UserInfo{}
	err := dbClient.Select(&userInfos, "sql/utils/select_user_profile.sql", &params)
	if err != nil {
		return map[string]UserInfo{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルの取得に失敗しました", Log: err.Error()}
	}
	res := map[string]UserInfo{}
	for _, v := range userInfos {
		res[v.UserId] = v
	}
	return res, nil
}

func CheckUserId(userIds []string, targetUserId string) bool {
	for _, userId := range userIds {
		if userId == targetUserId {
			return true
		}
	}
	return false
}

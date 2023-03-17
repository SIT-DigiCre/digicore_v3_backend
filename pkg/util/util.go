package util

import (
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
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
	err := dbClient.Select(&fileInfos, "sql/util/select_file.sql", &params)
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
	err := dbClient.Select(&userInfos, "sql/util/select_user_profile.sql", &params)
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

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

type FileId struct {
	FileId string `db:"file_id"`
}

type FileInfo struct {
	FileId string `db:"file_id"`
	Name   string `db:"name"`
}

func GetFileInfo(dbClient db.Client, fileIds []string) ([]FileInfo, *response.Error) {
	if len(fileIds) == 0 {
		return []FileInfo{}, nil
	}
	params := struct {
		FileIds []string `twowaysql:"fileIds"`
	}{
		FileIds: fileIds,
	}
	fileInfos := []FileInfo{}
	err := dbClient.Select(&fileInfos, "sql/util/select_file.sql", &params)
	if err != nil {
		return []FileInfo{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルの取得に失敗しました", Log: err.Error()}
	}
	return fileInfos, nil
}

type UserID struct {
	UserId string `db:"user_id"`
}

type UserInfo struct {
	IconUrl  string `db:"icon_url"`
	UserId   string `db:"user_id"`
	Username string `db:"username"`
}

func GetUserInfo(dbClient db.Client, userIds []string) ([]UserInfo, *response.Error) {
	if len(userIds) == 0 {
		return []UserInfo{}, nil
	}
	params := struct {
		UserIds []string `twowaysql:"userIds"`
	}{
		UserIds: userIds,
	}
	userInfos := []UserInfo{}
	err := dbClient.Select(&userInfos, "sql/util/select_user_profile.sql", &params)
	if err != nil {
		return []UserInfo{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ファイルの取得に失敗しました", Log: err.Error()}
	}
	return userInfos, nil
}

package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func IdFromStudentNumber(dbClient db.Client, studentNumber string) (string, *response.Error) {
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	user := []struct {
		Id     string `db:"id"`
		Active bool   `db:"active"`
	}{}
	err := dbClient.Select(&user, "sql/user/select_id_from_student_number.sql", &params)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	if len(user) == 0 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザーが存在しません", Log: "ユーザーが存在しません"}
	}
	if !user[0].Active {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "無効なアカウントです", Log: fmt.Sprintf("non active user login(%s)", user[0].Id)}
	}
	return user[0].Id, nil
}

type profile struct {
	UserId            string `db:"user_id"`
	StudentNumber     string `db:"student_number"`
	Username          string `db:"username"`
	SchoolGrade       int    `db:"school_grade"`
	IconUrl           string `db:"icon_url"`
	DiscordUserId     string `db:"discord_userid"`
	ActiveLimit       string `db:"active_limit"`
	ShortIntroduction string `db:"short_introduction"`
	IsAdmin           bool   `db:"is_admin"`
}

func GetUserProfileFromUserId(dbClient db.Client, userId string) (profile, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	profiles := []profile{}
	err := dbClient.Select(&profiles, "sql/user/select_user_profile_from_user_id.sql", &params)
	if err != nil {
		return profile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(profiles) == 0 {
		return profile{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "プロフィールが有りません", Log: "no rows in result"}
	}
	return profiles[0], nil
}

type introduction struct {
	Introduction string `db:"introduction"`
}

func GetStudentNumbersFromUserIds(dbClient db.Client, userIds []string) (map[string]string, *response.Error) {
	if len(userIds) == 0 {
		return map[string]string{}, nil
	}
	params := struct {
		UserIds []string `twowaysql:"userIds"`
	}{
		UserIds: userIds,
	}
	rows := []struct {
		UserId        string         `db:"user_id"`
		StudentNumber sql.NullString `db:"student_number"`
	}{}
	err := dbClient.Select(&rows, "sql/user/select_student_numbers_from_user_ids.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "ユーザー情報の取得に失敗しました", Log: err.Error()}
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		if row.StudentNumber.Valid {
			result[row.UserId] = row.StudentNumber.String
		}
	}
	return result, nil
}

func GetUserIntroductionFromUserId(dbClient db.Client, userId string) (introduction, *response.Error) {
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	introductions := []introduction{}
	err := dbClient.Select(&introductions, "sql/user/select_user_introduction_from_user_id.sql", &params)
	if err != nil {
		return introduction{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if len(introductions) == 0 {
		return introduction{}, &response.Error{Code: http.StatusNotFound, Level: "Info", Message: "自己紹介が有りません", Log: "no rows in result"}
	}
	return introductions[0], nil
}

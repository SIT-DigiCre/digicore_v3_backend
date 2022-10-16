package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func IDFromStudentNumber(dbClient db.Client, studentNumber string) (string, *response.Error) {
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	user := []struct {
		ID string `db:"id"`
	}{}
	err := dbClient.Select(&user, "sql/user/select_id_from_student_number.sql", &params)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	if len(user) != 1 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: "一意制約に違反しているデータがあります"}
	}
	return user[0].ID, nil
}

type profile struct {
	UserID            string `db:"user_id"`
	StudentNumber     string `db:"student_number"`
	Username          string `db:"username"`
	SchoolGrade       int    `db:"school_grade"`
	IconUrl           string `db:"icon_url"`
	DiscordUserID     string `db:"discord_userid"`
	ActiveLimit       string `db:"active_limit"`
	ShortIntroduction string `db:"short_introduction"`
}

func GetUserProfileFromUserID(dbClient db.Client, userID string) (profile, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
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

func GetUserIntroductionFromUserID(dbClient db.Client, userID string) (introduction, *response.Error) {
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
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

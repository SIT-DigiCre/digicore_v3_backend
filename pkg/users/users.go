package users

import (
	"database/sql"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func IDFromStudentNumber(studentNumber string, db db.DBClient) (string, *response.Error) {
	query, err := db.Query.ReadFile("sql/users/select_id_from_student_number.sql")
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	params := struct {
		StudentNumber string `twowaysql:"studentNumber"`
	}{
		StudentNumber: studentNumber,
	}
	users := []struct {
		ID string `db:"id"`
	}{}
	err = db.Client.Select(&users, string(query), &params)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	if len(users) != 1 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: "一意制約に違反しているデータがあります"}
	}
	return users[0].ID, nil
}

type Profile struct {
	UserId                string `db:"user_id"`
	StudentNumber         string `db:"student_number"`
	Username              string `db:"username"`
	SchoolGrade           int    `db:"school_grade"`
	IconUrl               string `db:"icon_url"`
	DiscordUserID         string `db:"discord_userid"`
	ActiveLimit           string `db:"active_limit"`
	ShortSelfIntroduction string `db:"short_self_introduction"`
}

func GetUserProfileFromUserID(userID string, db db.DBClient) (Profile, *response.Error) {
	query, err := db.Query.ReadFile("sql/users/select_user_profile_from_user_id.sql")
	if err != nil {
		return Profile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	params := struct {
		UserID string `twowaysql:"userID"`
	}{
		UserID: userID,
	}
	profile := []Profile{}
	err = db.Client.Select(&profile, string(query), &params)
	if err != sql.ErrNoRows {
		return Profile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "プロフィールが有りません", Log: err.Error()}
	}
	if err != nil {
		return Profile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: err.Error()}
	}
	if len(profile) != 1 {
		return Profile{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "DBエラーが発生しました", Log: "一意制約に違反しているデータがあります"}
	}
	return profile[0], nil
}
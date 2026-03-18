package user

import (
	"time"
    "github.com/labstack/echo/v4"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

type UserProfileLink struct {
	UrlLink string `db:"url_link"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func GetUserProfileLinksFromUserId(dbClient db.Client, userId string) ([]UserProfileLink, *response.Error) {
	if userId == "" {
		return []UserProfileLink{}, nil
	}
	params := struct {
		UserId string `twowaysql:"user_id"`
	}{
		UserId: userId,
	}
	linkUrls := []UserProfileLink{}
	err := dbClient.Select(&linkUrls, "sql/user/select_user_profile_links_from_user_id.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return linkUrls, nil
}

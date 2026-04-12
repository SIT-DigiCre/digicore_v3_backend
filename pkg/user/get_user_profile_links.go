package user

import (
	"time"
	"net/http"
    "github.com/getkin/kin-openapi/openapi_types"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

type UserProfileLink struct {
	Id        openapi_types.UUID    `db:"id"`
	LinkUrl   string    `db:"link_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func GetUserProfileLinksFromUserId(dbClient db.Client, userId string) ([]UserProfileLink, *response.Error) {
	if userId == "" {
		return []UserProfileLink{}, nil
	}
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	linkUrls := []UserProfileLink{}
	err := dbClient.Select(&linkUrls, "sql/user/select_user_profile_links_from_user_id.sql", &params)
	if err != nil {
		return []UserProfileLink{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "プロフィールリンクの取得に失敗しました", Log: err.Error()}
	}
	return linkUrls, nil
}

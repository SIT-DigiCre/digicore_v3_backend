package user

import (

	"github.com/labstack/echo/v4"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
)

func PostUserProfileLinks(c echo.Context, userID string, reqBody api.PostUserProfileLinksJSONRequestBody) error {
	client := c.Get("db").(*db.Client)

	params := struct {
		UserID  string `twowaysql:"user_id"`
		LinkUrl string `twowaysql:"link_url"`
	}{
		UserID:  userID,
		LinkUrl: reqBody.LinkUrl,
	}

	err := client.Exec("user/insert_user_profile_links", params)
	if err != nil {
		return response.InternalServerError(c, "プロフィールリンクの作成に失敗しました", err)
	}

	return response.Created(c, map[string]string{
		"user_id":  userID,
		"link_url": reqBody.LinkUrl,
	})
}

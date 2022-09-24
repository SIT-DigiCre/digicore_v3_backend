package users

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/labstack/echo/v4"
)

func GetUserMe(ctx echo.Context, dbClient db.Client) (api.ResGetUserMe, *response.Error) {
	userID := ctx.Get("user_id").(string)
	profile, err := GetUserProfileFromUserID(userID, dbClient)
	if err != nil {
		return api.ResGetUserMe{}, err
	}
	myProfile := api.ResGetUserMe{
		ActiveLimit:           profile.ActiveLimit,
		DiscordUserid:         profile.DiscordUserID,
		IconUrl:               profile.IconUrl,
		SchoolGrade:           profile.SchoolGrade,
		ShortSelfIntroduction: profile.ShortSelfIntroduction,
		StudentNumber:         profile.StudentNumber,
		UserId:                profile.UserId,
		Username:              profile.Username,
	}
	return myProfile, nil
}

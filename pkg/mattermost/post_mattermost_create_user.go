package mattermost

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/mattermost/mattermost-server/v5/model"
)

func PostMattermostCreateUser(ctx echo.Context, dbClient db.Client, requestBody api.ReqPostMattermostCreateuser) (api.ResPostMattermostCreateuser, *response.Error) {
	res := api.ResPostMattermostCreateuser{}
	userId := ctx.Get("user_id").(string)
	profile, err := user.GetUserProfileFromUserId(dbClient, userId)
	if err != nil {
		return api.ResPostMattermostCreateuser{}, err
	}
	email := fmt.Sprintf("%s@shibaura-it.ac.jp", profile.StudentNumber)

	client := model.NewAPIv4Client(env.MattermostURL)
	client.Login(env.MattermostAdminAccount, env.MattermostAdminPassword)
	team, rerr := client.RegenerateTeamInviteId(env.MattermostDigicreTeamID)
	if rerr.Error != nil {
		return api.ResPostMattermostCreateuser{}, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "招待IDの生成に失敗しました", Log: rerr.Error.Error()}
	}
	inviteID := team.InviteId

	user := &model.User{
		Username: requestBody.Username,
		Email:    email,
		Password: requestBody.Password,
		Nickname: requestBody.Nickname,
	}
	createdUser, rerr := client.CreateUserWithInviteId(user, inviteID)
	if rerr.Error != nil {
		return api.ResPostMattermostCreateuser{}, &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "ユーザーの追加に失敗しました", Log: rerr.Error.Error()}
	}
	res.Username = createdUser.Username

	if 0 < len(profile.IconUrl) {
		hres, err := http.Get(profile.IconUrl)
		if err != nil {
			return res, nil
		}
		defer hres.Body.Close()
		iconData, err := ioutil.ReadAll(hres.Body)
		if err != nil {
			return res, nil
		}
		client.SetProfileImage(createdUser.Id, iconData)
	}

	return res, nil
}

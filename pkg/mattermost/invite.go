package mattermost

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
	"github.com/mattermost/mattermost-server/v5/model"
)

var inviteID string
var inviteIDGeneratedAt time.Time

type ResponseInviteLink struct {
	URL string `json:"url"`
}
type ResponseError struct {
	Message string `json:"message"`
}

func getURLFromInviteID(inviteID string) string {
	return fmt.Sprintf("%s/signup_user_complete/?id=%s", env.MattermostURL, inviteID)
}

// Get Mattermost invite link
// @Router /mattermost/invite
// @Security Authorization
// @Success 200 {object}  ResponseInviteLink
// @Failure 500 {object} ResponseError
func GetInviteURL(e echo.Context) error {
	if len(inviteID) == 0 || inviteIDGeneratedAt.Unix() < time.Now().Unix() - 60 * 60 * 24 {
		client := model.NewAPIv4Client(env.MattermostURL)
		client.Login(env.MattermostAdminAccount, env.MattermostAdminPassword)

		team, _ := client.RegenerateTeamInviteId(env.MattermostTeamID)
		if team == nil {
			return e.JSON(http.StatusInternalServerError, ResponseError{ Message: "招待リンクの生成に失敗しました。" })
		}
		inviteID = team.InviteId
		inviteIDGeneratedAt = time.Now()
	}
	URL := getURLFromInviteID(inviteID)
	return e.JSON(http.StatusOK, ResponseInviteLink{ URL: URL })
}

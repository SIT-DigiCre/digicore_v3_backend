package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/discord"
	"github.com/labstack/echo/v4"
)

type RequestUpdateDiscordId struct {
	Code string `json:"code"`
}

type ResponseUpdateDiscordId struct {
	Error string `json:"error"`
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscordId body RequestUpdateDiscordId true "discord oauth code"
// @Security Authorization
// @Router /user/my/discord [put]
// @Success 200 {object} ResponseUpdateDiscordId
func (c Context) UpdateDiscordId(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscordId{Error: err.Error()})
	}
	request := RequestUpdateDiscordId{}
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscordId{Error: "データの読み込みに失敗しました"})
	}
	accessToken, err := discord.GetAccessToken(request.Code)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscordId{Error: err.Error()})
	}
	dicordID, err := discord.GetID(accessToken)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscordId{Error: err.Error()})
	}
	_, err = c.DB.Exec(`UPDATE UserProfile SET discord_userid = ? WHERE user_id = ?`,
		dicordID, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscordId{Error: "更新に失敗しました"})
	}
	return e.JSON(http.StatusOK, ResponseUpdateDiscordId{})
}

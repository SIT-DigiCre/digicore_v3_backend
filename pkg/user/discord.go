package user

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/discord"
	"github.com/labstack/echo/v4"
)

type RequestUpdateDiscord struct {
	Code string `json:"code"`
}

type ResponseUpdateDiscord struct {
	Error string `json:"error"`
}

// Update discord id
// @Accept json
// @Param RequestUpdateDiscord body RequestUpdateDiscord true "discord oauth code"
// @Router /user/my/disocrd [put]
// @Success 200 {object} ResponseUpdateDiscord
func (c Context) UpdateDiscordId(e echo.Context) error {
	userId, err := GetUserId(&e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscord{})
	}
	request := RequestUpdateDiscord{}
	if err := e.Bind(&request); err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscord{Error: err.Error()})
	}
	accessToken, err := discord.GetAccessToken(request.Code)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscord{Error: err.Error()})
	}
	dicordID, err := discord.GetID(accessToken)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseUpdateDiscord{Error: err.Error()})
	}
	_, err = c.DB.Exec(`UPDATE UserProfile SET discord_userid = ? WHERE user_id = UUID_TO_BIN(?)`,
		dicordID, userId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseSetMyPrivateProfile{Error: err.Error()})
	}
	return e.JSON(http.StatusOK, ResponseSetMyPrivateProfile{})
}

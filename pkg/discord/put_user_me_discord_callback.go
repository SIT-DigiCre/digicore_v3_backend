package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/labstack/echo/v4"
)

func PutUserMeDiscordCallback(ctx echo.Context, dbClient db.TransactionClient, requestBody api.ReqPutUserMeDiscordCallback) (api.ResGetUserMe, *response.Error) {
	userId := ctx.Get("user_id").(string)

	accessToken, err := getAccessToken(requestBody.Code)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	dicordUserId, err := getUserIdfromToken(accessToken)
	if err != nil {
		return api.ResGetUserMe{}, err
	}
	err = updateUserDiscordUserId(dbClient, userId, dicordUserId)
	if err != nil {
		return api.ResGetUserMe{}, err
	}

	return user.GetUserMe(ctx, dbClient)
}

func getAccessToken(code string) (string, *response.Error) {
	req, err := http.NewRequest("POST", "https://discordapp.com/api/oauth2/token", strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&code=%s&redirect_uri=%s", env.DiscordClientId, env.DiscordClientSecret, code, loginRedirectUrl)))
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&accessToken)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if res.StatusCode != 200 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: fmt.Sprintf("return %d not 200 (%s)", res.StatusCode, body)}
	}
	return accessToken.AccessToken, nil
}

func getUserIdfromToken(token string) (string, *response.Error) {
	req, err := http.NewRequest("GET", "https://discordapp.com/api/users/@me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	id := struct {
		Id string `json:"id"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&id)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	if res.StatusCode != 200 {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: fmt.Sprintf("return %d not 200", res.StatusCode)}
	}
	return id.Id, nil
}

func updateUserDiscordUserId(dbClient db.TransactionClient, userId string, discordUserId string) *response.Error {
	params := struct {
		UserId        string `twowaysql:"userId"`
		DiscordUserId string `twowaysql:"discordUserId"`
	}{
		UserId:        userId,
		DiscordUserId: discordUserId,
	}
	_, err := dbClient.Exec("sql/discord/update_user_discord_id.sql", &params, false)
	if err != nil {
		return &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "不明なエラーが発生しました", Log: err.Error()}
	}
	return nil
}

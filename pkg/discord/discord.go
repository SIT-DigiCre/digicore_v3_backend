package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
)

type Context struct {
}

func CreateContext() (Context, error) {
	context := Context{}
	return context, nil
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type IDResponse struct {
	ID string `json:"id"`
}

func GetAccessToken(code string) (string, error) {
	req, err := http.NewRequest("POST", "https://discordapp.com/api/oauth2/token", strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=authorization_code&code=%s&redirect_uri=%s", env.DiscordClientID, env.DiscordClientSecret, code, env.BackendRootURL+"/discord/oauth/callback")))
	if err != nil {
		return "", fmt.Errorf("アクセストークンの取得リクエストの生成に失敗しました")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("アクセストークンの取得リクエストに失敗しました")
	}
	accessToken := AccessTokenResponse{}
	err = json.NewDecoder(res.Body).Decode(&accessToken)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("アクセストークンの取得に失敗しました")
	}
	return accessToken.AccessToken, nil
}

func GetID(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://discordapp.com/api/users/@me", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	if err != nil {
		return "", fmt.Errorf("ユーザー情報の取得リクエストの生成に失敗しました")
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ユーザー情報の取得リクエストに失敗しました")
	}
	id := IDResponse{}
	err = json.NewDecoder(res.Body).Decode(&id)
	if err != nil || res.StatusCode != 200 {
		return "", fmt.Errorf("ユーザー情報の取得に失敗しました")
	}
	return id.ID, nil
}

package google

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	echo_session "github.com/ipfans/echo-session"
	"github.com/labstack/echo/v4"
)

type ResponseOAuthCallback struct {
	AccessToken string `json:"access_token"`
}

type UserInfoResponse struct {
	Email string `json:"email"`
}

const emailSuffix = "@shibaura-it.ac.jp"

func (u UserInfoResponse) validate() error {
	if !strings.HasSuffix(u.Email, emailSuffix) {
		return fmt.Errorf("suffix is not %s", emailSuffix)
	}
	return nil
}

// OAuth callback destination
// @Accept json
// @Router /google/oauth/callback [get]
// @Param code query string true "auth token"
// @Success 200 {object} ResponseOAuthCallback
func (c Context) OAuthCallback(e echo.Context) error {
	code := e.QueryParam("code")
	studentNumber, err := c.CheckGooleAccount(code)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontRootURL+"/login?")
	}
	userUuid, err := c.GetUserUuid(studentNumber)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontRootURL+"/login?")
	}
	sessionId, err := GetSessionId(&e, userUuid)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontRootURL+"/login?")
	}
	return e.Redirect(http.StatusFound, env.FrontRootURL+"/login?session="+sessionId)
}

func (c Context) CheckGooleAccount(code string) (string, error) {
	ctx := context.Background()
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("")
	}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?access_token="+token.AccessToken, nil)
	if err != nil {
		return "", fmt.Errorf("")
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("")
	}
	userInfo := UserInfoResponse{}
	json.NewDecoder(res.Body).Decode(&userInfo)
	if err := userInfo.validate(); err != nil {
		return "", fmt.Errorf("")
	}
	return strings.TrimSuffix(userInfo.Email, emailSuffix), nil
}

func (c Context) GetUserUuid(studentNumber string) (string, error) {
	userUuid := ""
	if err := c.DB.QueryRow("SELECT id FROM User WHERE student_number = ?", studentNumber).Scan(&userUuid); err == sql.ErrNoRows {
		_, err := c.DB.Exec("INSERT INTO User (student_number) VALUES (?)", studentNumber)
		if err != nil {
			return "", err
		}
		if err := c.DB.QueryRow("SELECT id FROM User WHERE student_number = ?", studentNumber).Scan(&userUuid); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}
	return userUuid, nil
}

func GetSessionId(e *echo.Context, userUuid string) (string, error) {
	session := echo_session.Default(*e)
	session.Set("id", userUuid)
	session.Set("login", true)
	session.Save()
	sessionId, err := (*e).Cookie("session")
	if err != nil {
		return "", err
	}
	return sessionId.Value, nil
}

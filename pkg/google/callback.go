package google

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

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

// OAuth callback destination (login only)
// @Router /google/oauth/callback/login [get]
// @Param code query string true "oauth code"
// @Success 302 "send authorization code to frontend"
// @Header 302 {string}  Location "/logined?session={}"
func (c Context) OAuthCallbackLogin(e echo.Context) error {
	code := e.QueryParam("code")
	redirectURL := oauth2.SetAuthURLParam("redirect_uri", env.BackendRootURL+"/google/oauth/callback/login")
	studentNumber, err := c.CheckGooleAccount(code, redirectURL)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/logined?")
	}
	userId, err := c.GetUserId(studentNumber)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/logined?")
	}
	sessionId, err := GetJWT(userId)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/logined?")
	}
	return e.Redirect(http.StatusFound, env.FrontendRootURL+"/logined?session="+sessionId)
}

// OAuth callback destination (register only)
// @Router /google/oauth/callback/register [get]
// @Param code query string true "oauth code"
// @Success 302 "send authorization code to frontend"
// @Header 302 {string}  Location "/registered?session={}"
func (c Context) OAuthCallbackRegister(e echo.Context) error {
	code := e.QueryParam("code")
	redirectURL := oauth2.SetAuthURLParam("redirect_uri", env.BackendRootURL+"/google/oauth/callback/register")
	studentNumber, err := c.CheckGooleAccount(code, redirectURL)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/registered?")
	}
	userId, err := c.RegisterUser(studentNumber)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/registered?")
	}
	sessionId, err := GetJWT(userId)
	if err != nil {
		return e.Redirect(http.StatusFound, env.FrontendRootURL+"/registered?")
	}
	return e.Redirect(http.StatusFound, env.FrontendRootURL+"/registered?session="+sessionId)
}

func (c Context) CheckGooleAccount(code string, redirectURL oauth2.AuthCodeOption) (string, error) {
	ctx := context.Background()
	token, err := c.Config.Exchange(ctx, code, redirectURL)
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

func (c Context) RegisterUser(studentNumber string) (string, error) {
	_, err := c.DB.Exec("INSERT INTO User (student_number) VALUES (?)", studentNumber)
	if err != nil {
		return "", err
	}
	userId, err := c.GetUserId(studentNumber)
	if err != nil {
		return "", nil
	}
	if err := user.CreateDefault(c.DB, userId, studentNumber); err != nil {
		return "", err
	}
	return userId, nil
}

func (c Context) GetUserId(studentNumber string) (string, error) {
	userId := ""
	if err := c.DB.QueryRow("SELECT BIN_TO_UUID(id) FROM User WHERE student_number = ?", studentNumber).Scan(&userId); err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}
	return userId, nil
}

func GetJWT(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userId
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString([]byte(env.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

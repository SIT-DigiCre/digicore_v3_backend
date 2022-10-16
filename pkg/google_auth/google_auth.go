package google_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var gcpConfig *oauth2.Config
var signupUrl string
var signupRedirectUrl = env.FrontendRootURL + env.SignupRedirectPath
var loginUrl string
var loginRedirectUrl = env.FrontendRootURL + env.LoginRedirectPath

func init() {
	gcpSecretJson, _ := os.ReadFile("./config/gcp_secret.json")
	if gcpSecretJson == nil {
		logrus.Fatal("Not found ./config/gcp_secret.json")
	}
	gcpConfig, _ = google.ConfigFromJSON(gcpSecretJson, "https://www.googleapis.com/auth/userinfo.email")
	signupRedirectURL := oauth2.SetAuthURLParam("redirect_uri", signupRedirectUrl)
	signupUrl = gcpConfig.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce, signupRedirectURL)
	loginRedirectURL := oauth2.SetAuthURLParam("redirect_uri", loginRedirectUrl)
	loginUrl = gcpConfig.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce, loginRedirectURL)
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

func getStudentNumberfromGoogle(code string, redirectURL string) (string, *response.Error) {
	redirectParam := oauth2.SetAuthURLParam("redirect_uri", redirectURL)
	token, err := gcpConfig.Exchange(context.Background(), code, redirectParam)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?access_token=%s", token.AccessToken), nil)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	userInfo := UserInfoResponse{}
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "認証でエラーが発生しました", Log: err.Error()}
	}
	if err := userInfo.validate(); err != nil {
		return "", &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "使用出来ないアカウントです", Log: err.Error()}
	}
	studentNumber := strings.TrimSuffix(userInfo.Email, emailSuffix)
	return studentNumber, nil
}

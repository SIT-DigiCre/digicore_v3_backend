package google_auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var gcpConfig *oauth2.Config
var signupUrl string
var signupRedirectUrl = env.FrontendRootURL + "/signup/callback"
var loginUrl string
var loginRedirectUrl = env.FrontendRootURL + "/login/callback"

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

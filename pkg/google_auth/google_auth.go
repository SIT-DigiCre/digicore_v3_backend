package google_auth

import (
	"os"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var gcpConfig *oauth2.Config
var signupUrl string

func init() {
	gcpSecretJson, _ := os.ReadFile("./config/gcp_secret.json")
	if gcpSecretJson == nil {
		logrus.Fatal("Not found ./config/gcp_secret.json")
	}
	gcpConfig, _ = google.ConfigFromJSON(gcpSecretJson, "https://www.googleapis.com/auth/userinfo.email")
	signupRedirectURL := oauth2.SetAuthURLParam("redirect_uri", env.FrontendRootURL+"/signup/code")
	signupUrl = gcpConfig.AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.ApprovalForce, signupRedirectURL)
}

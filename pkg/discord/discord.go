package discord

import (
	"fmt"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
)

var loginUrl string
var loginRedirectUrl = env.FrontendRootURL + env.DiscordLoginRedirectPath

func init() {
	loginUrl = fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify", env.DiscordClientId, loginRedirectUrl)
}

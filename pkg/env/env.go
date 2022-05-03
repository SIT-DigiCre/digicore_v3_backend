package env

import "os"

var BackendRootURL = os.Getenv("BACKEND_ROOT_URL")

var FrontendRootURL = os.Getenv("FRONTEND_ROOT_URL")

var JWTSecret = os.Getenv("JWT_SECRET")
var DefaultIconURL = os.Getenv("DEFAULT_ICON_URL")

var DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
var DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
var DiscordRedirectURL = os.Getenv("DISCORD_REDIRECT_URL")

var SlackURL = os.Getenv("SLACK_URL")
var DiscordURL = os.Getenv("DISCORD_URL")

type Context struct {
}

func CreateContext() (Context, error) {
	context := Context{}
	return context, nil
}

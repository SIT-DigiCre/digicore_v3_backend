package env

import "os"

var FrontendRootURL = os.Getenv("FRONTEND_ROOT_URL")
var BackendRootURL = os.Getenv("BACKEND_ROOT_URL")

var SignupRedirectPath = os.Getenv("SIGNUP_REDIRECT_PATH")
var LoginRedirectPath = os.Getenv("LOGIN_REDIRECT_PATH")

var DiscordLoginRedirectPath = os.Getenv("DISCORD_LOGIN_REDIRECT_PATH")

var DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
var DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")

var Auth = os.Getenv("AUTH")

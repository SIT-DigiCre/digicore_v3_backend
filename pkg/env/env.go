package env

import "os"

var FrontRootURL = os.Getenv("FRONT_ROOT_URL")

var JWTSecret = os.Getenv("JWT_SECRET")
var DefaultIconURL = os.Getenv("DEFAULT_ICON_URL")

var DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
var DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
var DiscordRedirectURL = os.Getenv("DISCORD_REDIRECT_URL")

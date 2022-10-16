package env

import "os"

var FrontendRootURL = os.Getenv("FRONTEND_ROOT_URL")
var BackendRootURL = os.Getenv("BACKEND_ROOT_URL")

var DiscordClientID = os.Getenv("DISCORD_CLIENT_ID")
var DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")

var Auth = os.Getenv("AUTH")

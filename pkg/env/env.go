package env

import "os"

var FrontendRootUrl = os.Getenv("FRONTEND_ROOT_URL")
var BackendRootUrl = os.Getenv("BACKEND_ROOT_URL")

var SignupRedirectPath = os.Getenv("SIGNUP_REDIRECT_PATH")
var LoginRedirectPath = os.Getenv("LOGIN_REDIRECT_PATH")

var DiscordLoginRedirectPath = os.Getenv("DISCORD_LOGIN_REDIRECT_PATH")

var DefaultIconUrl = os.Getenv("DEFAULT_ICON_URL")

var DiscordClientId = os.Getenv("DISCORD_CLIENT_ID")
var DiscordClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")

var DiscordJoinUrl = os.Getenv("DISCORD_JOIN_URL")

var Auth = os.Getenv("AUTH")

var WasabiAccessKey = os.Getenv("WASABI_ACCESS_KEY")
var WasabiSecretKey = os.Getenv("WASABI_SECRET_KEY")
var WasabiEndpoint = "s3.ap-northeast-1-ntt.wasabisys.com"
var WasabiRegion = "ap-northeast-1"
var WasabiPrivateBucket = os.Getenv("WASABI_PRIVATE_BUCKET")
var WasabiPublicBucket = os.Getenv("WASABI_PUBLIC_BUCKET")
var WasabiDirectURLDomain = "s3.ap-northeast-1.wasabisys.com"

var MattermostURL = "https://mm.digicre.net"
var MattermostDigicreTeamID = os.Getenv("MATTERMOST_DIGICRE_TEAM_ID")
var MattermostAdminAccount = os.Getenv("MATTERMOST_ADMIN_ACCOUNT")
var MattermostAdminPassword = os.Getenv("MATTERMOST_ADMIN_PASSWORD")

var AdminGroup = os.Getenv("ADMIN_GROUP")

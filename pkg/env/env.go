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

var ConohaIdentityServerURL = os.Getenv("CONOHA_IDENTITY_SERVER_URL")
var ConohaObjectStorageServerURL = os.Getenv("CONOHA_OBJECT_STORAGE_SERVER_URL")

var ConohaTenantID = os.Getenv("CONOHA_TENANT_ID")
var ConohaAPIUserName = os.Getenv("CONOHA_API_USER_NAME")
var ConohaAPIUserPassword = os.Getenv("CONOHA_API_USER_PASSWORD")
var ConohaStorageContainerName = "corev3"
var ConohaFileUploadMaxSize = 104857600

var WasabiAccessKey = os.Getenv("WASABI_ACCESS_KEY")
var WasabiSecretKey = os.Getenv("WASABI_SECRET_KEY")
var WasabiEndpoint = "s3.ap-northeast-1-ntt.wasabisys.com"
var WasabiRegion = "ap-northeast-1"
var WasabiPrivateBucket = os.Getenv("WASABI_PRIVATE_BUCKET")
var WasabiPublicBucket = os.Getenv("WASABI_PUBLIC_BUCKET")
var WasabiDirectURLDomain = "s3.ap-northeast-1.wasabisys.com"

var MattermostURL = "https://mm.digicre.net"
var MattermostTeamID = os.Getenv("MATTERMOST_TEAM_ID")
var MattermostAdminAccount = os.Getenv("MATTERMOST_ADMIN_ACCOUNT")
var MattermostAdminPassword = os.Getenv("MATTERMOST_ADMIN_PASSWORD")

type Context struct {
}

func CreateContext() (Context, error) {
	context := Context{}
	return context, nil
}

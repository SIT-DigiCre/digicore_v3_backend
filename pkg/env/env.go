package env

import "os"

// 未設定チェックの対象として参照した環境変数名
var checkedNames []string

// Get は環境変数を取得し、未設定チェックの対象として記録します
// パッケージ初期化時に呼び出す想定のため、実行時に呼び出さないでください
func Get(name string) string {
	checkedNames = append(checkedNames, name)
	return os.Getenv(name)
}

// Missing は未設定（または空文字）の環境変数名を返します
func Missing() []string {
	var missing []string
	for _, name := range checkedNames {
		if os.Getenv(name) == "" {
			missing = append(missing, name)
		}
	}
	return missing
}

var FrontendRootUrl = Get("FRONTEND_ROOT_URL")
var BackendRootUrl = Get("BACKEND_ROOT_URL")

var SignupRedirectPath = Get("SIGNUP_REDIRECT_PATH")
var LoginRedirectPath = Get("LOGIN_REDIRECT_PATH")

var DiscordLoginRedirectPath = Get("DISCORD_LOGIN_REDIRECT_PATH")

var DefaultIconUrl = Get("DEFAULT_ICON_URL")

var DiscordClientId = Get("DISCORD_CLIENT_ID")
var DiscordClientSecret = Get("DISCORD_CLIENT_SECRET")

var DiscordJoinUrl = Get("DISCORD_JOIN_URL")

var Auth = Get("AUTH")

var WasabiAccessKey = Get("WASABI_ACCESS_KEY")
var WasabiSecretKey = Get("WASABI_SECRET_KEY")
var WasabiEndpoint = "s3.ap-northeast-1-ntt.wasabisys.com"
var WasabiRegion = "ap-northeast-1"
var WasabiPrivateBucket = Get("WASABI_PRIVATE_BUCKET")
var WasabiPublicBucket = Get("WASABI_PUBLIC_BUCKET")
var WasabiDirectURLDomain = "s3.ap-northeast-1.wasabisys.com"

var MattermostURL = "https://mm.digicre.net"
var MattermostDigicreTeamID = Get("MATTERMOST_DIGICRE_TEAM_ID")
var MattermostAdminAccount = Get("MATTERMOST_ADMIN_ACCOUNT")
var MattermostAdminPassword = Get("MATTERMOST_ADMIN_PASSWORD")
var MattermostWebHookURL = Get("MATTERMOST_WEB_HOOK_URL")

var ClubRoomStatusToken = Get("CLUB_ROOM_STATUS_TOKEN")

var ResendApiKey = Get("RESEND_API_KEY")
var AdminEmail = Get("ADMIN_EMAIL")

package env

import "os"

var FrontRootURL = os.Getenv("FRONT_ROOT_URL")

var JWTSecret = os.Getenv("JWT_SECRET")
var DefaultIconURL = os.Getenv("DEFAULT_ICON_URL")

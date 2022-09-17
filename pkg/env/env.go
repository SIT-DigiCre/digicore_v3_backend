package env

import "os"

var FrontendRootURL = os.Getenv("FRONTEND_ROOT_URL")
var BackendRootURL = os.Getenv("BACKEND_ROOT_URL")

var Auth = os.Getenv("AUTH")

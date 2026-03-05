package user

import "github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/authenticator"

func init() {
	authenticator.RegisterNonMemberAllowedPath("/user/me/reentry")
}

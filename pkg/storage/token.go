package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
)

// Request structs
type passwordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type auth struct {
	PasswordCredentials passwordCredentials `json:"passwordCredentials"`
	TenantID            string              `json:"tenantId"`
}

type tokenRequest struct {
	Auth auth `json:"auth"`
}

// Response structs
type token struct {
	ID       string `json:"id"`
	IssuedAt string `json:"issued_at"`
	Expires  string `json:"expires"`
}
type access struct {
	Token token `json:"token"`
}
type tokenResponse struct {
	Access access `json:"access"`
}

func GetToken() (string, error) {
	credential := passwordCredentials{Username: env.ConohaAPIUserName, Password: env.ConohaAPIUserPassword}
	auth := auth{PasswordCredentials: credential, TenantID: env.ConohaTenantID}
	tokenReq := tokenRequest{Auth: auth}
	reqBodyBytes, err := json.Marshal(tokenReq)
	if err != nil {
		return "", errors.New("オブジェクトストレージのトークン取得エラーです")
	}
	reqBody := string(reqBodyBytes)
	reqBodyIo := strings.NewReader(reqBody)
	resBodyBytes, err := httpRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/tokens",
			env.ConohaIdentityServerURL,
		),
		reqBodyIo,
		nil,
	)
	if err != nil {
		return "", errors.New("オブジェクトストレージのトークン取得エラーです")
	}
	tokenRes := &tokenResponse{}
	if err := json.Unmarshal(resBodyBytes, &tokenRes); err != nil {
		return "", errors.New("オブジェクトストレージのトークン取得エラーです")
	}
	return tokenRes.Access.Token.ID, nil
}

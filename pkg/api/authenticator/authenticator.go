package authenticator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/group"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var key *rsa.PrivateKey

const PermissionsClaim = "perm"

func CreateAuthenticator() ([]echo.MiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	authenticator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			ErrorHandler: handler,
			Options: openapi3filter.Options{
				AuthenticationFunc: Authenticate,
			},
			Skipper: urlSkipper,
		})
	return []echo.MiddlewareFunc{authenticator, Login}, nil
}

func CreateDebugAuthenticator() ([]echo.MiddlewareFunc, error) {
	return []echo.MiddlewareFunc{Login}, nil
}

func init() {
	key = loadKey()
}

// loadKey は環境変数 JWT_PRIVATE_KEY から PEM 形式の RSA 秘密鍵を読み込みます。
// 未設定の場合は従来どおりランダムに生成します（開発環境用）。
func loadKey() *rsa.PrivateKey {
	pemText := os.Getenv("JWT_PRIVATE_KEY")
	if pemText == "" {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(fmt.Sprintf("JWT 鍵の生成に失敗しました: %v", err))
		}
		return key
	}

	// 環境変数経由では改行が失われる可能性があるため復元する
	pemText = strings.ReplaceAll(pemText, `\n`, "\n")

	block, _ := pem.Decode([]byte(pemText))
	if block == nil {
		panic("JWT_PRIVATE_KEY: PEM デコードに失敗しました")
	}

	// PKCS#8 形式（BEGIN PRIVATE KEY）のみ対応
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("JWT_PRIVATE_KEY: 秘密鍵の解析に失敗しました（PKCS#8 形式で設定してください）: %v", err))
	}
	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		panic("JWT_PRIVATE_KEY: 秘密鍵が RSA ではありません")
	}
	return rsaKey
}

func handler(c echo.Context, err *echo.HTTPError) error {
	message := err.Message.(string)
	if err.Code == http.StatusForbidden {
		res := response.Error{Code: http.StatusForbidden, Level: "Info", Message: "閲覧する権限がありません", Log: message}
		return response.ErrorResponse(c, &res)
	}
	if err.Code == http.StatusUnauthorized {
		res := response.Error{Code: http.StatusUnauthorized, Level: "Info", Message: "ログインされていません", Log: message}
		return response.ErrorResponse(c, &res)
	}
	res := response.Error{Code: err.Code, Level: "Info", Message: message, Log: message}
	return response.ErrorResponse(c, &res)
}

func CreateToken(user_id string) (string, *response.Error) {
	t := jwt.New()
	if err := t.Set(jwt.SubjectKey, user_id); err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "JWTトークンにユーザーIDを設定できませんでした", Log: err.Error()}
	}
	if err := t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*24*30).Unix()); err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "JWTトークンに有効期限を設定できませんでした", Log: err.Error()}
	}
	dbClient := db.Open()
	claims, err := group.GetClaimsFromUserId(&dbClient, user_id)
	if err != nil {
		return "", err
	}
	if err := t.Set(PermissionsClaim, claims); err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "JWTトークンに権限情報を設定できませんでした", Log: err.Error()}
	}

	signed, rerr := jwt.Sign(t, jwt.WithKey(jwa.RS256, key))
	if rerr != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "JWTトークンの署名に失敗しました", Log: rerr.Error()}
	}
	token := string(signed)
	return token, nil
}

func urlSkipper(c echo.Context) bool {
	return strings.HasPrefix(c.Path(), "/metrics") || strings.HasPrefix(c.Path(), "/mattermost/cmd")
}

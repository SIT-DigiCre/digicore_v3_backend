package authenticator

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
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
	key, _ = rsa.GenerateKey(rand.Reader, 2048)
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
	if err := t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*72).Unix()); err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "JWTトークンに有効期限を設定できませんでした", Log: err.Error()}
	}
	claims, err := GetClaims(user_id)
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

func GetClaims(userId string) ([]string, *response.Error) {
	dbClient := db.Open()

	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	claims := []claim{}
	err := dbClient.Select(&claims, "sql/group/select_claim_group_from_user_id.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "権限一覧の取得に失敗しました", Log: err.Error()}
	}
	claims_str := []string{}
	for _, claim := range claims {
		claims_str = append(claims_str, claim.Claim)
	}
	return claims_str, nil
}

type claim struct {
	Claim string `db:"claim"`
}

func urlSkipper(c echo.Context) bool {
	// Skip OpenAPI request validation for metrics and mattermost command endpoints.
	// Additionally skip for POST /event because the event POST handler binds raw
	// JSON and performs its own validation/parsing (avoids validator tag issues
	// on generated time.Time fields).
	// Use Request.URL.Path instead of Echo's registered path to be more robust
	// against routing differences.
	// Check multiple path sources to be robust against routing/baseURL prefixes
	reqPath := c.Request().URL.Path
	rawPath := c.Request().URL.RawPath
	echoPath := c.Path()

	// helper to detect event path (matches /event or /.../event and optional trailing slash)
	isEventPath := func(p string) bool {
		if p == "" {
			return false
		}
		if p == "/event" || p == "/event/" {
			return true
		}
		if strings.HasSuffix(p, "/event") || strings.HasSuffix(p, "/event/") {
			return true
		}
		// also support paths like /api/v1/event
		parts := strings.Split(p, "/")
		for i := range parts {
			if parts[i] == "event" {
				return true
			}
		}
		return false
	}

	if c.Request().Method == http.MethodPost && (isEventPath(reqPath) || isEventPath(rawPath) || isEventPath(echoPath)) {
		return true
	}

	return strings.HasPrefix(reqPath, "/metrics") || strings.HasPrefix(reqPath, "/mattermost/cmd")
}

package authenticator

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
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
	t.Set(jwt.SubjectKey, user_id)
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*72).Unix())
	claims, err := GetClaims(user_id)
	if err != nil {
		return "", err
	}
	t.Set(PermissionsClaim, claims)

	signed, rerr := jwt.Sign(t, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "JWTの生成に失敗しました", Log: rerr.Error()}
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

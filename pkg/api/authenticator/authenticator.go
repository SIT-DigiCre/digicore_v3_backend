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
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
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
	res := response.Error{Code: http.StatusUnauthorized, Level: "Info", Message: message, Log: message}
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
	groups := []group{}
	err := dbClient.Select(&groups, "sql/group/select_claim_group_from_user_id.sql", &params)
	if err != nil {
		return nil, &response.Error{Code: http.StatusInternalServerError, Level: "Error", Message: "権限一覧の取得に失敗しました", Log: err.Error()}
	}
	claims := []string{}
	for _, group := range groups {
		groupId := convertClaimGroupID(group.GroupId)
		claims = append(claims, groupId)
	}
	return claims, nil
}

type group struct {
	GroupId string `db:"group_id"`
}

func convertClaimGroupID(groupId string) string {
	if groupId == env.AdminGroup {
		return "admin"
	}
	return groupId
}

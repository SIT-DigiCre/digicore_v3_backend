package validator

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var key *rsa.PrivateKey

func CreateValidator() ([]echo.MiddlewareFunc, error) {
	spec, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			ErrorHandler: handler,
			Options: openapi3filter.Options{
				AuthenticationFunc: Authenticate,
			},
		})
	return []echo.MiddlewareFunc{validator, Login}, nil
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

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		fmt.Printf("failed to sign token: %s", err)
		return "", &response.Error{Code: http.StatusInternalServerError, Level: "Info", Message: "JWTの生成に失敗しました", Log: err.Error()}
	}
	token := string(signed)
	return token, nil
}

package authenticator

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

var key *rsa.PrivateKey

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

func init() {
	key, _ = rsa.GenerateKey(rand.Reader, 2048)
}

func handler(c echo.Context, err *echo.HTTPError) error {
	message := err.Message.(string)
	res := response.Error{Code: http.StatusUnauthorized, Level: "Info", Message: message, Log: message}
	return response.ErrorResponse(c, &res)
}

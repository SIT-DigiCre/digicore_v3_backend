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
	t, err := CreateToken("23876862-33ef-11ed-a261-0242ac120002")
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}
	fmt.Printf("%s", string(t))

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

func CreateToken(user_id string) ([]byte, error) {
	t := jwt.New()
	t.Set(jwt.SubjectKey, user_id)
	t.Set(jwt.ExpirationKey, time.Now().Add(time.Hour*72).Unix())

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		fmt.Printf("failed to sign token: %s", err)
		return []byte{}, err
	}
	return signed, nil
}

package authenticator

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/deepmap/oapi-codegen/examples/petstore-expanded/echo/api"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func getJWT(input *openapi3filter.AuthenticationInput) string {
	token := input.RequestValidationInput.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

func validateJWT(token string) (jwt.Token, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, fmt.Errorf("security scheme")
	}
	return verifiedToken, nil
}

func getClaims(t jwt.Token) ([]string, error) {
	rawPerms, found := t.Get("perm")
	if !found {
		return make([]string, 0), nil
	}

	rawList, ok := rawPerms.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' claim is unexpected type'", "perm")
	}

	claims := make([]string, len(rawList))

	for i, rawClaim := range rawList {
		claims[i], ok = rawClaim.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d] is not a string", "perm", i)
		}
	}
	return claims, nil
}

func checkClaims(expectedClaims []string, t jwt.Token) error {
	claims, err := getClaims(t)
	if err != nil {
		return fmt.Errorf("getting claims from token: %w", err)
	}
	// Put the claims into a map, for quick access.
	claimsMap := make(map[string]bool, len(claims))
	for _, c := range claims {
		claimsMap[c] = true
	}

	for _, e := range expectedClaims {
		if !claimsMap[e] {
			return fmt.Errorf("security scheme")
		}
	}
	return nil
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}
	tokenString := getJWT(input)
	token, err := validateJWT(tokenString)
	if err != nil {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	err = checkClaims(input.Scopes, token)
	if err != nil {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}
	return nil
}

func Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			return next(c)
		}
		token, err := jwt.Parse([]byte(tokenString), jwt.WithVerify(false))
		if err != nil {
			return c.JSON(http.StatusBadRequest, api.Error{Message: err.Error()})
		}
		subjectKey, _ := token.Get(jwt.SubjectKey)
		c.Set("user_id", subjectKey)
		return next(c)
	}
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

package authenticator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return echo.ErrUnauthorized
	}
	tokenString := getJWT(input)
	token, err := validateJWT(tokenString)
	if err != nil {
		return echo.ErrUnauthorized
	}

	err = checkClaims(input.Scopes, token)
	if err != nil {
		return echo.ErrForbidden
	}
	return nil
}

func Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" || urlSkipper(c) {
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

func validateJWT(token string) (jwt.Token, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, fmt.Errorf("security scheme")
	}
	return verifiedToken, nil
}

func getJWT(input *openapi3filter.AuthenticationInput) string {
	token := input.RequestValidationInput.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	return token
}

func getClaims(t jwt.Token) ([]string, error) {
	rawPerms, found := t.Get(PermissionsClaim)
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

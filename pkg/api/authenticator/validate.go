package authenticator

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/db"
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

		userId, ok := subjectKey.(string)
		if !ok {
			return response.ErrorResponse(c, &response.Error{
				Code:    http.StatusUnauthorized,
				Level:   "Info",
				Message: "ログインされていません",
				Log:     "subject claim is not string",
			})
		}
		isMember, rerr := checkUserIsMember(userId)
		if rerr != nil {
			return response.ErrorResponse(c, rerr)
		}
		if !isMember && !isNonMemberAllowedPath(c.Path()) {
			return response.ErrorResponse(c, &response.Error{
				Code:    http.StatusForbidden,
				Level:   "Info",
				Message: "無効なアカウントです",
				Log:     fmt.Sprintf("non member user access denied(%s, %s)", userId, c.Path()),
			})
		}
		return next(c)
	}
}

func checkUserIsMember(userId string) (bool, *response.Error) {
	dbClient := db.Open()
	params := struct {
		UserId string `twowaysql:"userId"`
	}{
		UserId: userId,
	}
	rows := []struct {
		IsMember bool `db:"is_member"`
	}{}
	err := dbClient.Select(&rows, "sql/user/select_user_is_member_from_user_id.sql", &params)
	if err != nil {
		return false, &response.Error{
			Code:    http.StatusInternalServerError,
			Level:   "Error",
			Message: "不明なエラーが発生しました",
			Log:     err.Error(),
		}
	}
	if len(rows) == 0 {
		return false, &response.Error{
			Code:    http.StatusUnauthorized,
			Level:   "Info",
			Message: "ログインされていません",
			Log:     fmt.Sprintf("user profile not found(%s)", userId),
		}
	}
	return rows[0].IsMember, nil
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
	fmt.Printf("%v %v\n", claims, expectedClaims)
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

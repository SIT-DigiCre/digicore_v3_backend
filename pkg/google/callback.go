package google

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseOAuthCallback struct {
}

// OAuth callback destination
// @Accept json
// @Router /google/oauth/callback [get]
// @Param code query string true "auth token"
// @Success 200 {object} ResponseOAuthCallback
func (c Context) OAuthCallback(e echo.Context) error {
	return e.JSON(http.StatusOK, ResponseOAuthCallback{})
}

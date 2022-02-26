package google

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type ResponseOAuthCallback struct {
	AccessToken string `json:"access_token"`
}

type UserInfoResponse struct {
	Email string `json:"email"`
}

const hostDomain = "shibaura-it.ac.jp"

func validateIdToken(idToken string, c *oauth2.Config) error {
	set, err := jwk.Fetch(context.Background(), "https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return fmt.Errorf("failed to parse JWK: %w", err)
	}
	token, err := jwt.Parse([]byte(idToken), jwt.WithKeySet(set))
	if err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}
	if token.Issuer() != "https://accounts.google.com" {
		return fmt.Errorf("iss not match")
	}
	if len(token.Audience()) != 1 || token.Audience()[0] != c.ClientID {
		return fmt.Errorf("audience is not %s", c.ClientID)
	}
	if token.Expiration().Unix() < time.Now().Unix() {
		return fmt.Errorf("expiration of term of validity: %s", token.Expiration())
	}

	hd, ok := token.Get("hd")
	if !ok || hd.(string) != hostDomain {
		return fmt.Errorf("host domain is not %s", hostDomain)
	}
	return nil
}

// OAuth callback destination
// @Accept json
// @Router /google/oauth/callback [get]
// @Param code query string true "auth token"
// @Success 200 {object} ResponseOAuthCallback
func (c Context) OAuthCallback(e echo.Context) error {
	code := e.QueryParam("code")

	ctx := context.Background()
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseOAuthCallback{})
	}
	idToken := token.Extra("id_token").(string)
	if validateIdToken(idToken, c.Config) != nil {
		return e.JSON(http.StatusBadRequest, ResponseOAuthCallback{})
	}

	return e.JSON(http.StatusOK, ResponseOAuthCallback{AccessToken: idToken})
}

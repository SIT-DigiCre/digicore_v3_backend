// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /event)
	GetEvent(ctx echo.Context, params GetEventParams) error

	// (GET /event/{eventID})
	GetEventEventID(ctx echo.Context, eventID string) error

	// (GET /event/{eventID}/{reservationID})
	GetEventEventIDReservationID(ctx echo.Context, eventID string, reservationID string) error

	// (DELETE /event/{eventID}/{reservationID}/me)
	DeleteEventEventIDReservationIDMe(ctx echo.Context, eventID string, reservationID string) error

	// (PUT /event/{eventID}/{reservationID}/me)
	PutEventEventIDReservationIDMe(ctx echo.Context, eventID string, reservationID string) error

	// (GET /login)
	GetLogin(ctx echo.Context) error

	// (POST /login/callback)
	PostLoginCallback(ctx echo.Context) error

	// (GET /signup)
	GetSignup(ctx echo.Context) error

	// (POST /signup/callback)
	PostSignupCallback(ctx echo.Context) error

	// (GET /status)
	GetStatus(ctx echo.Context) error

	// (GET /tool)
	GetTool(ctx echo.Context) error

	// (GET /user)
	GetUser(ctx echo.Context, params GetUserParams) error

	// (GET /user/me)
	GetUserMe(ctx echo.Context) error

	// (PUT /user/me)
	PutUserMe(ctx echo.Context) error

	// (GET /user/me/discord)
	GetUserMeDiscord(ctx echo.Context) error

	// (PUT /user/me/discord/callback)
	PutUserMeDiscordCallback(ctx echo.Context) error

	// (GET /user/me/introduction)
	GetUserMeIntroduction(ctx echo.Context) error

	// (PUT /user/me/introduction)
	PutUserMeIntroduction(ctx echo.Context) error

	// (GET /user/me/payment)
	GetUserMePayment(ctx echo.Context) error

	// (PUT /user/me/payment)
	PutUserMePayment(ctx echo.Context) error

	// (GET /user/me/private)
	GetUserMePrivate(ctx echo.Context) error

	// (PUT /user/me/private)
	PutUserMePrivate(ctx echo.Context) error

	// (GET /user/{userID})
	GetUserUserID(ctx echo.Context, userID string) error

	// (GET /user/{userID}/introduction)
	GetUserUserIDIntroduction(ctx echo.Context, userID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetEvent converts echo context to params.
func (w *ServerInterfaceWrapper) GetEvent(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetEventParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "seed" -------------

	err = runtime.BindQueryParameter("form", true, false, "seed", ctx.QueryParams(), &params.Seed)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter seed: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEvent(ctx, params)
	return err
}

// GetEventEventID converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventEventID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventID" -------------
	var eventID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventID", runtime.ParamLocationPath, ctx.Param("eventID"), &eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEventEventID(ctx, eventID)
	return err
}

// GetEventEventIDReservationID converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventEventIDReservationID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventID" -------------
	var eventID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventID", runtime.ParamLocationPath, ctx.Param("eventID"), &eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventID: %s", err))
	}

	// ------------- Path parameter "reservationID" -------------
	var reservationID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationID", runtime.ParamLocationPath, ctx.Param("reservationID"), &reservationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEventEventIDReservationID(ctx, eventID, reservationID)
	return err
}

// DeleteEventEventIDReservationIDMe converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteEventEventIDReservationIDMe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventID" -------------
	var eventID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventID", runtime.ParamLocationPath, ctx.Param("eventID"), &eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventID: %s", err))
	}

	// ------------- Path parameter "reservationID" -------------
	var reservationID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationID", runtime.ParamLocationPath, ctx.Param("reservationID"), &reservationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteEventEventIDReservationIDMe(ctx, eventID, reservationID)
	return err
}

// PutEventEventIDReservationIDMe converts echo context to params.
func (w *ServerInterfaceWrapper) PutEventEventIDReservationIDMe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventID" -------------
	var eventID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventID", runtime.ParamLocationPath, ctx.Param("eventID"), &eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventID: %s", err))
	}

	// ------------- Path parameter "reservationID" -------------
	var reservationID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationID", runtime.ParamLocationPath, ctx.Param("reservationID"), &reservationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutEventEventIDReservationIDMe(ctx, eventID, reservationID)
	return err
}

// GetLogin converts echo context to params.
func (w *ServerInterfaceWrapper) GetLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetLogin(ctx)
	return err
}

// PostLoginCallback converts echo context to params.
func (w *ServerInterfaceWrapper) PostLoginCallback(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostLoginCallback(ctx)
	return err
}

// GetSignup converts echo context to params.
func (w *ServerInterfaceWrapper) GetSignup(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetSignup(ctx)
	return err
}

// PostSignupCallback converts echo context to params.
func (w *ServerInterfaceWrapper) PostSignupCallback(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostSignupCallback(ctx)
	return err
}

// GetStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetStatus(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStatus(ctx)
	return err
}

// GetTool converts echo context to params.
func (w *ServerInterfaceWrapper) GetTool(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTool(ctx)
	return err
}

// GetUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetUser(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUserParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "seed" -------------

	err = runtime.BindQueryParameter("form", true, false, "seed", ctx.QueryParams(), &params.Seed)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter seed: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUser(ctx, params)
	return err
}

// GetUserMe converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserMe(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserMe(ctx)
	return err
}

// PutUserMe converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMe(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMe(ctx)
	return err
}

// GetUserMeDiscord converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserMeDiscord(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserMeDiscord(ctx)
	return err
}

// PutUserMeDiscordCallback converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMeDiscordCallback(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMeDiscordCallback(ctx)
	return err
}

// GetUserMeIntroduction converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserMeIntroduction(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserMeIntroduction(ctx)
	return err
}

// PutUserMeIntroduction converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMeIntroduction(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMeIntroduction(ctx)
	return err
}

// GetUserMePayment converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserMePayment(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserMePayment(ctx)
	return err
}

// PutUserMePayment converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMePayment(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMePayment(ctx)
	return err
}

// GetUserMePrivate converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserMePrivate(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserMePrivate(ctx)
	return err
}

// PutUserMePrivate converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMePrivate(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMePrivate(ctx)
	return err
}

// GetUserUserID converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserUserID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "userID", runtime.ParamLocationPath, ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserUserID(ctx, userID)
	return err
}

// GetUserUserIDIntroduction converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserUserIDIntroduction(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userID" -------------
	var userID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "userID", runtime.ParamLocationPath, ctx.Param("userID"), &userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserUserIDIntroduction(ctx, userID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/event", wrapper.GetEvent)
	router.GET(baseURL+"/event/:eventID", wrapper.GetEventEventID)
	router.GET(baseURL+"/event/:eventID/:reservationID", wrapper.GetEventEventIDReservationID)
	router.DELETE(baseURL+"/event/:eventID/:reservationID/me", wrapper.DeleteEventEventIDReservationIDMe)
	router.PUT(baseURL+"/event/:eventID/:reservationID/me", wrapper.PutEventEventIDReservationIDMe)
	router.GET(baseURL+"/login", wrapper.GetLogin)
	router.POST(baseURL+"/login/callback", wrapper.PostLoginCallback)
	router.GET(baseURL+"/signup", wrapper.GetSignup)
	router.POST(baseURL+"/signup/callback", wrapper.PostSignupCallback)
	router.GET(baseURL+"/status", wrapper.GetStatus)
	router.GET(baseURL+"/tool", wrapper.GetTool)
	router.GET(baseURL+"/user", wrapper.GetUser)
	router.GET(baseURL+"/user/me", wrapper.GetUserMe)
	router.PUT(baseURL+"/user/me", wrapper.PutUserMe)
	router.GET(baseURL+"/user/me/discord", wrapper.GetUserMeDiscord)
	router.PUT(baseURL+"/user/me/discord/callback", wrapper.PutUserMeDiscordCallback)
	router.GET(baseURL+"/user/me/introduction", wrapper.GetUserMeIntroduction)
	router.PUT(baseURL+"/user/me/introduction", wrapper.PutUserMeIntroduction)
	router.GET(baseURL+"/user/me/payment", wrapper.GetUserMePayment)
	router.PUT(baseURL+"/user/me/payment", wrapper.PutUserMePayment)
	router.GET(baseURL+"/user/me/private", wrapper.GetUserMePrivate)
	router.PUT(baseURL+"/user/me/private", wrapper.PutUserMePrivate)
	router.GET(baseURL+"/user/:userID", wrapper.GetUserUserID)
	router.GET(baseURL+"/user/:userID/introduction", wrapper.GetUserUserIDIntroduction)

}

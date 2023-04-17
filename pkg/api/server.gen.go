// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /blog/blog)
	GetBlogBlog(ctx echo.Context, params GetBlogBlogParams) error

	// (POST /blog/blog)
	PostBlogBlog(ctx echo.Context) error

	// (DELETE /blog/blog/{blogId})
	DeleteBlogBlogBlogId(ctx echo.Context, blogId string) error

	// (GET /blog/blog/{blogId})
	GetBlogBlogBlogId(ctx echo.Context, blogId string) error

	// (PUT /blog/blog/{blogId})
	PutBlogBlogBlogId(ctx echo.Context, blogId string) error

	// (GET /blog/tag)
	GetBlogTag(ctx echo.Context, params GetBlogTagParams) error

	// (POST /blog/tag)
	PostBlogTag(ctx echo.Context) error

	// (DELETE /blog/tag/{tagId})
	DeleteBlogTagTagId(ctx echo.Context, tagId string) error

	// (GET /blog/tag/{tagId})
	GetBlogTagTagId(ctx echo.Context, tagId string) error

	// (PUT /blog/tag/{tagId})
	PutBlogTagTagId(ctx echo.Context, tagId string) error

	// (GET /event)
	GetEvent(ctx echo.Context, params GetEventParams) error

	// (GET /event/{eventId})
	GetEventEventId(ctx echo.Context, eventId string) error

	// (GET /event/{eventId}/{reservationId})
	GetEventEventIdReservationId(ctx echo.Context, eventId string, reservationId string) error

	// (DELETE /event/{eventId}/{reservationId}/me)
	DeleteEventEventIdReservationIdMe(ctx echo.Context, eventId string, reservationId string) error

	// (PUT /event/{eventId}/{reservationId}/me)
	PutEventEventIdReservationIdMe(ctx echo.Context, eventId string, reservationId string) error

	// (GET /group)
	GetGroup(ctx echo.Context, params GetGroupParams) error

	// (GET /group/{groupId})
	GetGroupGroupId(ctx echo.Context, groupId string) error

	// (GET /login)
	GetLogin(ctx echo.Context) error

	// (POST /login/callback)
	PostLoginCallback(ctx echo.Context) error

	// (POST /mattermost/create_user)
	PostMattermostCreateUser(ctx echo.Context) error

	// (GET /payment)
	GetPayment(ctx echo.Context, params GetPaymentParams) error

	// (GET /payment/{paymentId})
	GetPaymentPaymentId(ctx echo.Context, paymentId string) error

	// (PUT /payment/{paymentId})
	PutPaymentPaymentId(ctx echo.Context, paymentId string) error

	// (GET /signup)
	GetSignup(ctx echo.Context) error

	// (POST /signup/callback)
	PostSignupCallback(ctx echo.Context) error

	// (GET /status)
	GetStatus(ctx echo.Context) error

	// (GET /status/club_room)
	GetStatusClubRoom(ctx echo.Context) error

	// (PUT /status/club_room)
	PutStatusClubRoom(ctx echo.Context) error

	// (GET /storage/myfile)
	GetStorageMyfile(ctx echo.Context) error

	// (POST /storage/myfile)
	PostStorageMyfile(ctx echo.Context) error

	// (GET /storage/{fileId})
	GetStorageFileId(ctx echo.Context, fileId string) error

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

	// (PUT /user/me/renewal)
	PutUserMeRenewal(ctx echo.Context) error

	// (GET /user/{userId})
	GetUserUserId(ctx echo.Context, userId string) error

	// (GET /user/{userId}/introduction)
	GetUserUserIdIntroduction(ctx echo.Context, userId string) error

	// (GET /work/tag)
	GetWorkTag(ctx echo.Context, params GetWorkTagParams) error

	// (POST /work/tag)
	PostWorkTag(ctx echo.Context) error

	// (DELETE /work/tag/{tagId})
	DeleteWorkTagTagId(ctx echo.Context, tagId string) error

	// (GET /work/tag/{tagId})
	GetWorkTagTagId(ctx echo.Context, tagId string) error

	// (PUT /work/tag/{tagId})
	PutWorkTagTagId(ctx echo.Context, tagId string) error

	// (GET /work/work)
	GetWorkWork(ctx echo.Context, params GetWorkWorkParams) error

	// (POST /work/work)
	PostWorkWork(ctx echo.Context) error

	// (DELETE /work/work/{workId})
	DeleteWorkWorkWorkId(ctx echo.Context, workId string) error

	// (GET /work/work/{workId})
	GetWorkWorkWorkId(ctx echo.Context, workId string) error

	// (PUT /work/work/{workId})
	PutWorkWorkWorkId(ctx echo.Context, workId string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetBlogBlog converts echo context to params.
func (w *ServerInterfaceWrapper) GetBlogBlog(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBlogBlogParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "authorId" -------------

	err = runtime.BindQueryParameter("form", true, false, "authorId", ctx.QueryParams(), &params.AuthorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBlogBlog(ctx, params)
	return err
}

// PostBlogBlog converts echo context to params.
func (w *ServerInterfaceWrapper) PostBlogBlog(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostBlogBlog(ctx)
	return err
}

// DeleteBlogBlogBlogId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteBlogBlogBlogId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "blogId" -------------
	var blogId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "blogId", runtime.ParamLocationPath, ctx.Param("blogId"), &blogId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter blogId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteBlogBlogBlogId(ctx, blogId)
	return err
}

// GetBlogBlogBlogId converts echo context to params.
func (w *ServerInterfaceWrapper) GetBlogBlogBlogId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "blogId" -------------
	var blogId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "blogId", runtime.ParamLocationPath, ctx.Param("blogId"), &blogId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter blogId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBlogBlogBlogId(ctx, blogId)
	return err
}

// PutBlogBlogBlogId converts echo context to params.
func (w *ServerInterfaceWrapper) PutBlogBlogBlogId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "blogId" -------------
	var blogId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "blogId", runtime.ParamLocationPath, ctx.Param("blogId"), &blogId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter blogId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutBlogBlogBlogId(ctx, blogId)
	return err
}

// GetBlogTag converts echo context to params.
func (w *ServerInterfaceWrapper) GetBlogTag(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetBlogTagParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBlogTag(ctx, params)
	return err
}

// PostBlogTag converts echo context to params.
func (w *ServerInterfaceWrapper) PostBlogTag(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostBlogTag(ctx)
	return err
}

// DeleteBlogTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteBlogTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteBlogTagTagId(ctx, tagId)
	return err
}

// GetBlogTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) GetBlogTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBlogTagTagId(ctx, tagId)
	return err
}

// PutBlogTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) PutBlogTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutBlogTagTagId(ctx, tagId)
	return err
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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEvent(ctx, params)
	return err
}

// GetEventEventId converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventEventId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventId" -------------
	var eventId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventId", runtime.ParamLocationPath, ctx.Param("eventId"), &eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEventEventId(ctx, eventId)
	return err
}

// GetEventEventIdReservationId converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventEventIdReservationId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventId" -------------
	var eventId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventId", runtime.ParamLocationPath, ctx.Param("eventId"), &eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventId: %s", err))
	}

	// ------------- Path parameter "reservationId" -------------
	var reservationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationId", runtime.ParamLocationPath, ctx.Param("reservationId"), &reservationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetEventEventIdReservationId(ctx, eventId, reservationId)
	return err
}

// DeleteEventEventIdReservationIdMe converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteEventEventIdReservationIdMe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventId" -------------
	var eventId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventId", runtime.ParamLocationPath, ctx.Param("eventId"), &eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventId: %s", err))
	}

	// ------------- Path parameter "reservationId" -------------
	var reservationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationId", runtime.ParamLocationPath, ctx.Param("reservationId"), &reservationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteEventEventIdReservationIdMe(ctx, eventId, reservationId)
	return err
}

// PutEventEventIdReservationIdMe converts echo context to params.
func (w *ServerInterfaceWrapper) PutEventEventIdReservationIdMe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "eventId" -------------
	var eventId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "eventId", runtime.ParamLocationPath, ctx.Param("eventId"), &eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter eventId: %s", err))
	}

	// ------------- Path parameter "reservationId" -------------
	var reservationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "reservationId", runtime.ParamLocationPath, ctx.Param("reservationId"), &reservationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter reservationId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutEventEventIdReservationIdMe(ctx, eventId, reservationId)
	return err
}

// GetGroup converts echo context to params.
func (w *ServerInterfaceWrapper) GetGroup(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetGroupParams
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
	err = w.Handler.GetGroup(ctx, params)
	return err
}

// GetGroupGroupId converts echo context to params.
func (w *ServerInterfaceWrapper) GetGroupGroupId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "groupId" -------------
	var groupId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "groupId", runtime.ParamLocationPath, ctx.Param("groupId"), &groupId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter groupId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetGroupGroupId(ctx, groupId)
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

// PostMattermostCreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) PostMattermostCreateUser(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostMattermostCreateUser(ctx)
	return err
}

// GetPayment converts echo context to params.
func (w *ServerInterfaceWrapper) GetPayment(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{"account"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPaymentParams
	// ------------- Optional query parameter "year" -------------

	err = runtime.BindQueryParameter("form", true, false, "year", ctx.QueryParams(), &params.Year)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter year: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPayment(ctx, params)
	return err
}

// GetPaymentPaymentId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPaymentPaymentId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "paymentId" -------------
	var paymentId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "paymentId", runtime.ParamLocationPath, ctx.Param("paymentId"), &paymentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter paymentId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"account"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetPaymentPaymentId(ctx, paymentId)
	return err
}

// PutPaymentPaymentId converts echo context to params.
func (w *ServerInterfaceWrapper) PutPaymentPaymentId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "paymentId" -------------
	var paymentId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "paymentId", runtime.ParamLocationPath, ctx.Param("paymentId"), &paymentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter paymentId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{"account"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutPaymentPaymentId(ctx, paymentId)
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

// GetStatusClubRoom converts echo context to params.
func (w *ServerInterfaceWrapper) GetStatusClubRoom(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStatusClubRoom(ctx)
	return err
}

// PutStatusClubRoom converts echo context to params.
func (w *ServerInterfaceWrapper) PutStatusClubRoom(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutStatusClubRoom(ctx)
	return err
}

// GetStorageMyfile converts echo context to params.
func (w *ServerInterfaceWrapper) GetStorageMyfile(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStorageMyfile(ctx)
	return err
}

// PostStorageMyfile converts echo context to params.
func (w *ServerInterfaceWrapper) PostStorageMyfile(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostStorageMyfile(ctx)
	return err
}

// GetStorageFileId converts echo context to params.
func (w *ServerInterfaceWrapper) GetStorageFileId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "fileId" -------------
	var fileId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "fileId", runtime.ParamLocationPath, ctx.Param("fileId"), &fileId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fileId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetStorageFileId(ctx, fileId)
	return err
}

// GetTool converts echo context to params.
func (w *ServerInterfaceWrapper) GetTool(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

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

// PutUserMeRenewal converts echo context to params.
func (w *ServerInterfaceWrapper) PutUserMeRenewal(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutUserMeRenewal(ctx)
	return err
}

// GetUserUserId converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserUserId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserUserId(ctx, userId)
	return err
}

// GetUserUserIdIntroduction converts echo context to params.
func (w *ServerInterfaceWrapper) GetUserUserIdIntroduction(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUserUserIdIntroduction(ctx, userId)
	return err
}

// GetWorkTag converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkTag(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWorkTagParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkTag(ctx, params)
	return err
}

// PostWorkTag converts echo context to params.
func (w *ServerInterfaceWrapper) PostWorkTag(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostWorkTag(ctx)
	return err
}

// DeleteWorkTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteWorkTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteWorkTagTagId(ctx, tagId)
	return err
}

// GetWorkTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkTagTagId(ctx, tagId)
	return err
}

// PutWorkTagTagId converts echo context to params.
func (w *ServerInterfaceWrapper) PutWorkTagTagId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "tagId" -------------
	var tagId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "tagId", runtime.ParamLocationPath, ctx.Param("tagId"), &tagId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tagId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutWorkTagTagId(ctx, tagId)
	return err
}

// GetWorkWork converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkWork(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetWorkWorkParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "authorId" -------------

	err = runtime.BindQueryParameter("form", true, false, "authorId", ctx.QueryParams(), &params.AuthorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkWork(ctx, params)
	return err
}

// PostWorkWork converts echo context to params.
func (w *ServerInterfaceWrapper) PostWorkWork(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostWorkWork(ctx)
	return err
}

// DeleteWorkWorkWorkId converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteWorkWorkWorkId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "workId" -------------
	var workId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "workId", runtime.ParamLocationPath, ctx.Param("workId"), &workId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter workId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteWorkWorkWorkId(ctx, workId)
	return err
}

// GetWorkWorkWorkId converts echo context to params.
func (w *ServerInterfaceWrapper) GetWorkWorkWorkId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "workId" -------------
	var workId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "workId", runtime.ParamLocationPath, ctx.Param("workId"), &workId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter workId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWorkWorkWorkId(ctx, workId)
	return err
}

// PutWorkWorkWorkId converts echo context to params.
func (w *ServerInterfaceWrapper) PutWorkWorkWorkId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "workId" -------------
	var workId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "workId", runtime.ParamLocationPath, ctx.Param("workId"), &workId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter workId: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutWorkWorkWorkId(ctx, workId)
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

	router.GET(baseURL+"/blog/blog", wrapper.GetBlogBlog)
	router.POST(baseURL+"/blog/blog", wrapper.PostBlogBlog)
	router.DELETE(baseURL+"/blog/blog/:blogId", wrapper.DeleteBlogBlogBlogId)
	router.GET(baseURL+"/blog/blog/:blogId", wrapper.GetBlogBlogBlogId)
	router.PUT(baseURL+"/blog/blog/:blogId", wrapper.PutBlogBlogBlogId)
	router.GET(baseURL+"/blog/tag", wrapper.GetBlogTag)
	router.POST(baseURL+"/blog/tag", wrapper.PostBlogTag)
	router.DELETE(baseURL+"/blog/tag/:tagId", wrapper.DeleteBlogTagTagId)
	router.GET(baseURL+"/blog/tag/:tagId", wrapper.GetBlogTagTagId)
	router.PUT(baseURL+"/blog/tag/:tagId", wrapper.PutBlogTagTagId)
	router.GET(baseURL+"/event", wrapper.GetEvent)
	router.GET(baseURL+"/event/:eventId", wrapper.GetEventEventId)
	router.GET(baseURL+"/event/:eventId/:reservationId", wrapper.GetEventEventIdReservationId)
	router.DELETE(baseURL+"/event/:eventId/:reservationId/me", wrapper.DeleteEventEventIdReservationIdMe)
	router.PUT(baseURL+"/event/:eventId/:reservationId/me", wrapper.PutEventEventIdReservationIdMe)
	router.GET(baseURL+"/group", wrapper.GetGroup)
	router.GET(baseURL+"/group/:groupId", wrapper.GetGroupGroupId)
	router.GET(baseURL+"/login", wrapper.GetLogin)
	router.POST(baseURL+"/login/callback", wrapper.PostLoginCallback)
	router.POST(baseURL+"/mattermost/create_user", wrapper.PostMattermostCreateUser)
	router.GET(baseURL+"/payment", wrapper.GetPayment)
	router.GET(baseURL+"/payment/:paymentId", wrapper.GetPaymentPaymentId)
	router.PUT(baseURL+"/payment/:paymentId", wrapper.PutPaymentPaymentId)
	router.GET(baseURL+"/signup", wrapper.GetSignup)
	router.POST(baseURL+"/signup/callback", wrapper.PostSignupCallback)
	router.GET(baseURL+"/status", wrapper.GetStatus)
	router.GET(baseURL+"/status/club_room", wrapper.GetStatusClubRoom)
	router.PUT(baseURL+"/status/club_room", wrapper.PutStatusClubRoom)
	router.GET(baseURL+"/storage/myfile", wrapper.GetStorageMyfile)
	router.POST(baseURL+"/storage/myfile", wrapper.PostStorageMyfile)
	router.GET(baseURL+"/storage/:fileId", wrapper.GetStorageFileId)
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
	router.PUT(baseURL+"/user/me/renewal", wrapper.PutUserMeRenewal)
	router.GET(baseURL+"/user/:userId", wrapper.GetUserUserId)
	router.GET(baseURL+"/user/:userId/introduction", wrapper.GetUserUserIdIntroduction)
	router.GET(baseURL+"/work/tag", wrapper.GetWorkTag)
	router.POST(baseURL+"/work/tag", wrapper.PostWorkTag)
	router.DELETE(baseURL+"/work/tag/:tagId", wrapper.DeleteWorkTagTagId)
	router.GET(baseURL+"/work/tag/:tagId", wrapper.GetWorkTagTagId)
	router.PUT(baseURL+"/work/tag/:tagId", wrapper.PutWorkTagTagId)
	router.GET(baseURL+"/work/work", wrapper.GetWorkWork)
	router.POST(baseURL+"/work/work", wrapper.PostWorkWork)
	router.DELETE(baseURL+"/work/work/:workId", wrapper.DeleteWorkWorkWorkId)
	router.GET(baseURL+"/work/work/:workId", wrapper.GetWorkWorkWorkId)
	router.PUT(baseURL+"/work/work/:workId", wrapper.PutWorkWorkWorkId)

}

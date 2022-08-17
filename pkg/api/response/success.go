package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SuccessResponse(ctx echo.Context, body any) error {
	return ctx.JSON(http.StatusOK, body)
}

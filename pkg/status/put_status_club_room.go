package status

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func PutStatusClubRoom(ctx echo.Context, requestBody api.ReqPutStatusClubRoom) (api.ResGetStatusClubRoom, *response.Error) {
	clubRoomLock = requestBody.Lock
	return GetStatusClubRoom(ctx)
}

package status

import (
	"net/http"

	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/env"
	"github.com/labstack/echo/v4"
)

func PutStatusClubRoom(ctx echo.Context, requestBody api.ReqPutStatusClubRoom) (api.ResGetStatusClubRoom, *response.Error) {

	if requestBody.Token != env.ClubRoomStatusToken {
		return api.ResGetStatusClubRoom{}, &response.Error{Code: http.StatusBadRequest, Level: "Info", Message: "編集する権限がありません", Log: requestBody.Token + ":" + env.ClubRoomStatusToken + "no edit permission"}
	}

	clubRoomLock = requestBody.Lock
	return GetStatusClubRoom(ctx)
}

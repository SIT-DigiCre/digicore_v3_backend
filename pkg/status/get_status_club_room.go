package status

import (
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api"
	"github.com/SIT-DigiCre/digicore_v3_backend/pkg/api/response"
	"github.com/labstack/echo/v4"
)

func GetStatusClubRoom(ctx echo.Context) (api.ResGetStatusClubRoom, *response.Error) {
	return api.ResGetStatusClubRoom{Lock: clubRoomLock}, nil
}
